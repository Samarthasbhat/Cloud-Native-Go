package main

import (
	"fmt"
	// "html"
	"net/http"
	"sync"
	"text/template"
	"time"
)

func Debounce (fn func(string), wait time.Duration) func(string){
	var mu sync.Mutex
	var timer *time.Timer
	var lastInput string

	return func(input string){
		mu.Lock()
		defer mu.Unlock()

		lastInput = input

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, func(){
			fn(lastInput)
		})
	}
}

var tpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html> 
<head>
	<title> Go Search </title>
</head>
<body> 
	<h2>Search Example</h2>
	<form action="/search" method="GET">
		<input type="text" name="q" placeholder="Type something..." >
		<button type="submit">Search</button>
	</form>
</body>
</html>`))


func main(){
	// Actual search function
	 search := func (query string) {
		fmt.Println("Searching for:", query)
	 }

	 debounceSearch := Debounce(search, 2*time.Second)

	 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		 tpl.Execute(w, nil)
	 })

	 http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		if query == "" {
			// http.Error(w, "No query provided", http.StatusBadRequest)

			fmt.Fprintf(w, 
				`<!Doctype html>
			<html>
				<head>
					<title>Error</title>
				</head>
				<body>
					<h2>No query provided</h2>
					<a href="/">Go Back</a>
				</body>
			</html>`)
		}else{
		fmt.Fprintf(w, "Search request received for query: %s", query)
		}
		debounceSearch(query)
		
	 })
	 fmt.Println("Server started at :8080")
	 http.ListenAndServe(":8080", nil)
}