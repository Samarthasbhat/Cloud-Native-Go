package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// EventType defines the type of operation
type EventType byte

const (
	EventDelete EventType = iota + 1
	EventPut
)

// Event represents a transaction log entry
type Event struct {
	Sequence  uint
	EventType EventType
	Key       string
	Value     string
}

// TransactionLogger interface for logging
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}

// Database connection params
type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

// PostgresTransactionLogger struct
type PostgresTransactionLogger struct {
	events chan Event     // channel for sending events
	errors chan error     // channel for receiving errors
	db     *sql.DB        // database handle
}

// WritePut logs a PUT event
func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// WriteDelete logs a DELETE event
func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

// Err returns error channel
func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

// NewPostgresTransactionLogger initializes logger
func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		config.host, config.dbName, config.user, config.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

// Run starts listening for events and writing them to DB
func (l *PostgresTransactionLogger) Run() {
	l.events = make(chan Event, 16)
	l.errors = make(chan error, 1)

	go func() {
		query := `INSERT INTO transactions(event_type, key, value) VALUES($1, $2, $3)`

		for e := range l.events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				l.errors <- err
			}
		}
	}()
}

// ReadEvents replays the transaction log
func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var e Event
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}
			outEvent <- e
		}

		if err = rows.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

// verifyTableExists checks if transactions table exists
func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'transactions'
		)`
	var exists bool
	err := l.db.QueryRow(query).Scan(&exists)
	return exists, err
}

// createTable creates the transactions table if missing
func (l *PostgresTransactionLogger) createTable() error {
	query := `
		CREATE TABLE transactions (
			sequence SERIAL PRIMARY KEY,
			event_type INT,
			key TEXT,
			value TEXT
		)`
	_, err := l.db.Exec(query)
	return err
}

// Example usage
func main() {
	logger, err := NewPostgresTransactionLogger(PostgresDBParams{
		host:     "localhost",
		dbName:   "kvs",
		user:     "test",
		password: "hunter2",
	})
	if err != nil {
		panic(err)
	}

	logger.Run()

	// Write sample events
	logger.WritePut("foo", "bar")
	logger.WriteDelete("baz")

	// Read back events
	events, errs := logger.ReadEvents()
	for e := range events {
		fmt.Printf("Event: %+v\n", e)
	}
	for err := range errs {
		fmt.Printf("Error: %v\n", err)
	}
}
