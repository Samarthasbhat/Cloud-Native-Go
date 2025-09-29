package main

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq" // PostgreSQL driver
)

type EventType byte

const (
    EventDelete EventType = iota + 1
    EventPut
)

type Event struct {
    Sequence  uint
    EventType EventType
    Key       string
    Value     string
}

type TransactionLogger interface {
    WriteDelete(key string)
    WritePut(key, value string)
    Err() <-chan error
    ReadEvents() (<-chan Event, <-chan error)
    Run()
}

type PostgresDBParams struct {
    dbName   string
    host     string
    user     string
    password string
}

type PostgresTransactionLogger struct {
    events chan Event
    errors chan error
    db     *sql.DB
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
    l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
    l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
    return l.errors
}

func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
connStr := fmt.Sprintf(
    "host=%s port=5432 dbname=%s user=%s password=%s sslmode=disable",
    config.host, config.dbName, config.user, config.password,
)



    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to connect to db: %w", err)
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

func main() {
    // Use your actual host IP address here
    dbParams := PostgresDBParams{
        host: "192.168.1.5",     // Must match the IP allowed in pg_hba.conf
        dbName:   "kvs",             // Ensure this database exists
        user:     "test",            // Ensure this user exists and has access
        password: "hunter2",         // Use the correct password
    }

    logger, err := NewPostgresTransactionLogger(dbParams)
    if err != nil {
        fmt.Printf("Database setup failed: %v\n", err)
        return
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
