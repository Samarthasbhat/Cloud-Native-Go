package main

// These are the requirement for transaction log
// - Sequence number, Event type, Key, Value

import (
	// "fmt"
	"os"
	)

// transaction logger interface

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
}

type Event struct {
	Sequence uint64 // a unique record ID
	EventType EventType // The action taken
	Key string // The key affected by this transaction
	value string // The value of a PUT the transaction
}

type FileTransactionLogger struct{
	events chan <- Event
	errors <-chan error
	lastSequence uint64
	file *os.file
}

type EventType byte

const ( 
	_                     = iota
	EventDelete EventType = iota
	EventPut 
)

// Storing state in a transaction log file
// Pros: No downstream, Technically straightforward
// Cons: Harder to scale, Uncontrolled growth

func (l *FileTransactionLogger) WritePut(key, value string){
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func(l *FileTransactionLogger) WriteDelete(key string){
	l.events <-Event{EventType: EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error{
	return l.errors
}