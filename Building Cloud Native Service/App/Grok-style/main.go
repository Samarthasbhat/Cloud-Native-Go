package main

import (
	"fmt"
	"net/http"
	"strings"
)

var grokAnswers = map[string]string{
	"go":       "Go: Google's programming language that's faster than your coffee kicking in.",
    "channel":  "Channels: like walkie-talkies for your goroutines.",
    "concurrency": "Concurrency: letting your code juggle tasks like a caffeinated squirrel ",
}
func GoHome(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html") 
    fmt.Fprintln(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Welcome to Grok</title>
            <style>
                body {
                    background-color: #0d1117;
                    color: #f0f6fc;
                    font-family: Arial, sans-serif;
                    text-align: center;
                    padding-top: 100px;
                }
                h1 {
                    font-size: 3em;
                    color: #58a6ff;
                }
                p {
                    font-size: 1.2em;
                }
                .search-box {
                    margin-top: 20px;
                }
                input[type="text"] {
                    width: 300px;
                    padding: 12px 20px;
                    border-radius: 8px;
                    border: 1px solid #30363d;
                    background-color: #161b22;
                    color: #f0f6fc;
                    font-size: 1em;
                    outline: none;
                    transition: all 0.3s ease;
                }
                input[type="text"]:focus {
                    border-color: #58a6ff;
                    box-shadow: 0 0 8px rgba(88,166,255,0.6);
                }
                button {
                    padding: 12px 20px;
                    margin-left: 10px;
                    border: none;
                    border-radius: 8px;
                    background-color: #58a6ff;
                    color: white;
                    font-size: 1em;
                    cursor: pointer;
                    transition: background 0.3s ease;
                }
                button:hover {
                    background-color: #1f6feb;
                }
            </style>
        </head>
        <body>
            <h1> Welcome to Go Search</h1>
            <p>Go-powered API is alive!</p>

            <div class="search-box">
                <form action="/gosearch" method="get">
                    <input type="text" name="question" placeholder="Ask Related to GO..." required>
                    <button type="submit">Search</button>
                </form>
            </div>
        </body>
        </html>
    `)
}


func SearchGoStyle(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    question := strings.ToLower(r.FormValue("question"))

    for key, answer := range grokAnswers {
        if strings.Contains(question, key) {
            w.Header().Set("Content-Type", "text/html")
            fmt.Fprintf(w, `
                <!DOCTYPE html>
                <html>
                <head>
                    <title>Grok Answer</title>
                    <style>
                        body {
                            background-color: #0d1117;
                            color: #f0f6fc;
                            font-family: Arial, sans-serif;
                            text-align: center;
                            padding-top: 100px;
                        }
                        h1 { color: #58a6ff; }
                        a { color: #58a6ff; text-decoration: none; }
                    </style>
                </head>
                <body>
                    <h1>Answer:</h1>
                    <p>%s</p>
                    <br>
                    <a href="/go">Ask another question</a>
                </body>
                </html>
            `, answer)
            return
        }
    }

    // No answer found
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>No Answer</title>
            <style>
                body {
                    background-color: #0d1117;
                    color: #f0f6fc;
                    font-family: Arial, sans-serif;
                    text-align: center;
                    padding-top: 100px;
                }
                a { color: #58a6ff; text-decoration: none; }
            </style>
        </head>
        <body>
            <h1>Sorry, no answer found.</h1>
            <a href="/go">Try again</a>
        </body>
        </html>
    `)
}


func main() {
	http.HandleFunc("/go", GoHome)
	http.HandleFunc("/gosearch", SearchGoStyle)
	fmt.Println("Starting server on :8080")	
	http.ListenAndServe(":8080", nil)
}