package main

import (
	"net/http"
)

// Building an HTTP server with net/http

func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello net/http!"))
}

func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	// Implementation of ListenAndServeTLS
	return nil
}

func main() {
	http.HandleFunc("/", helloGoHandler)
	// http.ListenAndServe(":8080", nil)
	err :=	http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	if err!= nil {
		panic(err)
	}
}
