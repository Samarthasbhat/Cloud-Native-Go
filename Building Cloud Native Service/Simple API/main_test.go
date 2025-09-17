package main


import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	const key = "create-Key"
	const value = "create-Value"

	var val interface{}
	var contains bool 

	defer delete(store, key)

	_,contains = store[key]
	if contains {
		t.Errorf("store should not contain key %q", key)
	}

	// Err should be nil
	err := Put(key,value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}


func TestGet(t *testing.T) {
	const key = "get-Key"
	const value = "get-Value"

	var val string
	var err error

	defer delete(store, key)

	// Read a non-thing
	val, err = Get(key)
	if err == nil{
		t.Error("expected error")
	}
	if !errors.Is(err, ErrorNoSuchKey){
		t.Error("Unexpected error:", err)
	}

	store[key] = value

	val, err = Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	const key = "delete-key"
	const value = "delete-value"

	var contains bool

	defer delete(store, key)

	store[key] = value

	_, contains = store[key]
	if !contains {
		t.Error("key/value doesn't exist")
	}

	Delete(key)

	_, contains = store[key]
	if contains {
		t.Error("Delete failed")
	}
}