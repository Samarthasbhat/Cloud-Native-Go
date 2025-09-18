package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

var (
	tasks []string
	mu    sync.Mutex
)

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>TODO App</title>
	<style>
		body {
			background-color: #282c34;
			color: Black;
			font-family: Arial, sans-serif;
		}
		.container {
			width: 350px;
			margin: 50px auto;
			padding: 20px;
			background: #ffffffff;
			border-radius: 10px;
		}
		input[type=text] {
			padding: 8px;
			width: 70%;
			margin-right: 10px;
		}
		input[type=submit] {
			padding: 8px 12px;
			cursor: pointer;
		}
		.delete {
			background-color: #e74c3c;
			border: none;
			padding: 8px 12px;
			color: white;
			}
		.add {
			background-color: #2ecc71;
			border: none;
			border-radius: 50%;
			padding: 8px 12px;
		}
		h1{
			text-align: center;
			color: #333;
			margin-bottom: 20px;
			opacity: 0.8;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>TODO List</h1>
		
		<!-- Add Task Form -->
		<form action="/add" method="post">
			<input type="text" name="task" placeholder="Enter new task" required>
			<input class="add" type="submit" value="+">
		</form>

		

		<!-- Task List with Delete -->
		<form action="/delete" method="post">
			<h2>Tasks</h2>
			<ul>
				{{range $i, $t := .}}
					
						<input type="checkbox" name="delete" value="{{$i}}">
						<label>{{$t}}</label>
					
				{{else}}
					<li>No tasks yet</li>
				{{end}}
			</ul>
			<input class="delete" type="submit" value="Delete">
		</form>
	</div>
</body>
</html>
`))

// Renders the main page
func Home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	tmpl.Execute(w, tasks)
}

// Add a task
func AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		task := r.FormValue("task")
		if task != "" {
			mu.Lock()
			tasks = append(tasks, task)
			mu.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete selected tasks
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		selected := r.Form["delete"]

		if len(selected) > 0 {
			mu.Lock()
			toDelete := make(map[int]bool)
			for _, v := range selected {
				idx, err := strconv.Atoi(v)
				if err == nil {
					toDelete[idx] = true
				}
			}

			var newTasks []string
			for i, task := range tasks {
				if !toDelete[i] { // keep only unchecked ones
					newTasks = append(newTasks, task)
				}
			}
			tasks = newTasks
			mu.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/add", AddTask)
	http.HandleFunc("/delete", DeleteTask)

	http.ListenAndServe(":8080", nil)
}
