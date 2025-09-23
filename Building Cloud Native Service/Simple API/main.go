package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/mux"
)

// This is creating the core component of the API Services with GET<PUT<DELETE functions

// var store = make(map[string]string)

var store = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key string, value string) error {
		// store[key] = value
		// return nil

		store.Lock()
		store.m[key] = value
		store.Unlock()

		return nil
	}

	var ErrorNoSuchKey = errors.New("no such key")

func Get(key string) (string, error) {
	// value, ok := store[key]

	// if !ok{
	// 	return "", ErrorNoSuchKey
	// }
	// return value, nil
	store.RLock()
		value, ok := store.m[key]
		store.RLock()

		if !ok{
			return " ", ErrorNoSuchKey
		}
		return value, nil
}

func Delete(key string) error{
	// delete(store, key)   //Inbuilt function to delete a key from a map
	// return nil
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)
	return nil
}

//Create Function

func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, 
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	err = Put(key, string(value))
	if err != nil {
		http.Error(w, 
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}


// Read Function

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //Retrieve "Key" from the request
	key := vars["key"]

	value, err := Get(key) // Get Value for key
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, 
			err.Error(),
			http.StatusNotFound)
		return
	}
	if err != nil{
		http.Error(w, 
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value)) //Write value to response
}

// Delete Handler
func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	errors := Delete(key)
	if errors != nil {
		http.Error(w, 
			errors.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main(){
r := mux.NewRouter()

	r.HandleFunc("/v1/keys/{key}", KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/keys/{key}", KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/keys/{key}", KeyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}



