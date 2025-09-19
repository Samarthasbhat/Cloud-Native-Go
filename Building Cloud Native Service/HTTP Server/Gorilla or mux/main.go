package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler for /
func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Gorilla/Mux!"))
}

// Handler for /products/{key}
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	fmt.Fprintf(w, "Product key: %s\n", key)
}

// Handler for /articles/{category}/
func ArticleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	fmt.Fprintf(w, "Articles in category: %s\n", category)
}

// Handler for /articles/{category}/{id:[0-9]+}
func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	id := vars["id"]
	fmt.Fprintf(w, "Article ID: %s in category: %s\n", id, category)
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", helloMuxHandler)
	r.HandleFunc("/products/{key}", ProductHandler)
	r.HandleFunc("/articles/{category}/", ArticleCategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	// Start server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
