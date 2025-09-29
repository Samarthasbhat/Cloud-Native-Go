package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
)

type EventType byte

const (
	_                     = iota // iota == 0; ignore this value
	EventDelete EventType = iota // iota == 1
	EventPut                     // iota == 2; implicitly repeat last
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type TransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error
	lastSequence uint64   // The last used event sequence number
	file         *os.File // The location of the transaction log
	wg           *sync.WaitGroup
}

func (l *TransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventPut, Key: key, Value: url.QueryEscape(value)}
}

func (l *TransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *TransactionLogger) Err() <-chan error {
	return l.errors
}

func NewTransactionLogger(filename string) (*TransactionLogger, error) {
	var err error
	var l TransactionLogger = TransactionLogger{wg: &sync.WaitGroup{}}

	// Open the transaction log file for reading and writing.
	l.file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &l, nil
}

func (l *TransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// Start retrieving events from the events channel and writing them
	// to the transaction log
	go func() {
		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- fmt.Errorf("cannot write to log file: %w", err)
			}

			l.wg.Done()
		}
	}()
}

func (l *TransactionLogger) Wait() {
	l.wg.Wait()
}

func (l *TransactionLogger) Close() error {
	l.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.file.Close()
}

func (l *TransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	_, err := l.file.Seek(0,0)
		if err != nil {
		errCh := make(chan error, 1)
		close(errCh)
		return nil, errCh
	}
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			if e.Sequence <= l.lastSequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			uv, err := url.QueryUnescape(e.Value)
			if err != nil {
				outError <- fmt.Errorf("value decoding failure: %w", err)
				return
			}
			e.Value = uv

			l.lastSequence = e.Sequence

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func main() {
	// Create new transaction logger
	logger, err := NewTransactionLogger("transaction.log")
	if err != nil{
		panic(err)
	}
	// Start the logger
	logger.Run()

	// write some test cases
	logger.WritePut("user1","ABC")
	logger.WritePut("user2","XYZ")
	logger.WriteDelete("user1")

	// Wait and Close
	logger.Close()

	//Reopen and reply events

	logger2, err := NewTransactionLogger("transaction.log")
	if err != nil {
		panic(err)
	}

	events, errors := logger2.ReadEvents()

	fmt.Println("Replying events from log: ")
	for{
		select{
		case e, ok := <-events:
			if !ok {
				events=nil
			}else {
				fmt.Printf("Seq=%d Type=%d Key=%s Value=%s\n",
					e.Sequence, e.EventType, e.Key, e.Value)
			}
		case err, ok := <-errors:
			if ok{
				fmt.Println("Error:", err)
			}
			errors = nil
		}
		if events == nil && errors == nil{
			break
		}
	}
}