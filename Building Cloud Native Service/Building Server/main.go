package main

import (
	"net/http"
)

// Building an HTTP server with net/http

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello net/http!"))
}

func main() {
	http.HandleFunc("/hello", helloGoHandler)
	http.ListenAndServe(":8080", nil)
}
