package main

import (
	"errors"
)

// This is creating the core component of the API Services with GET<PUT<DELETE functions

var store = make(map[string]string) 

func Put(key string, value string) error {
		store[key] = value
		return nil
	}

	var ErrorNoSuchKey = errors.New("no such key")

func Get(key string) (string, error) {
	value, ok := store[key]

	if !ok{
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func Delete(key string) error{
	delete(store, key)   //Inbuilt function to delete a key from a map
	return nil
}


