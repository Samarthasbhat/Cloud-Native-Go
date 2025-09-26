package main


import (
	"database/sql"
	_ "github.com/lib/pq"  // Import the driver package
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

