package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import the driver package
)

// Importing a database drivers

//Implementing PostgresTransactionLogger

type TransactionLogger interface{
	WriteDelete(key string)
	WritePut(key,value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

// Required for Database access

type PostgresDBParams struct{
	dbName 	 string
	host   	 string
	user     string
	password string
}

// variables   

logger, err = NewPostgresTransactionLogger(PostgresDBParams{
	host: "localhost",
	dbName: "kvs",
	user: "test",
	password: "hunter2"
})

type PostgresTransactionLogger struct{
	events chan<- Event  //Write-only channel for sending events
	errors <-chan error //Read-only channel for receiving errors
	db *sql.DB    // Database access interface
}

func (l *PostgresTransactionLogger) WritePut(key,value string){
		l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string){
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error{
	return l.errors
}



func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error){
	connStr := fmt.Sprintf("host=%s dbame=%s user=%s password=%s",
	config.host, config.dbName, config.user, config.password)

	db, err := sql.Open("postgres", connStr)
	if err!= nil{
		return  nil, fmt.Errorf("failed to open db: %w", err)
	}
	logger := &PostgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err!= nil {
		return nil, fmt.Errorf("failed tp verify table exists: %w", err)
	}
	if !exists{
		if err = logger.createTable(); err != nil {
			return  nil, fmt.Errorf("failed to create table: %w", err)

		}
	}
	return logger
}