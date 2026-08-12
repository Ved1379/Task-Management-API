package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var tasks []Task
var nextID = 1

func Homehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Task Management API")
}

func main() {
	fmt.Println("Server started on port 8080")
	http.HandleFunc("/", Homehandler)
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/tasks/", taskHandler)
	http.ListenAndServe(":8080", nil)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if r.URL.Path == "/tasks" {

		if r.Method == http.MethodGet {
			for _, task := range tasks {
				fmt.Fprintln(w, "ID:", task.ID)
				fmt.Fprintln(w, "Title:", task.Title)
				fmt.Fprintln(w, "Description:", task.Description)
			}
		}

		if r.Method == http.MethodPost {

			var task Task

			err := json.NewDecoder(r.Body).Decode(&task)

			if err != nil {
				fmt.Fprintln(w, "Invalid JSON")
				return
			}

			task.ID = nextID
			nextID++

			tasks = append(tasks, task)

			fmt.Fprintln(w, "Task created:", task.Title)
		}

		return
	}

	id, err := strconv.Atoi(parts[2])

	if err != nil {
		fmt.Fprintln(w, "Invalid task Id")
		return
	}

	if r.Method == http.MethodGet {

		found := false

		for _, task := range tasks {
			if id == task.ID {
				json.NewEncoder(w).Encode(task)
				found = true
				return
			}
		}

		if !found {
			fmt.Fprintln(w, "Task not found")
		}
	}

	if r.Method == http.MethodPut {

		var updatedTask Task

		err := json.NewDecoder(r.Body).Decode(&updatedTask)

		if err != nil {
			fmt.Fprintln(w, "Invalid JSON")
			return
		}

		for i, task := range tasks {
			if id == task.ID {

				updatedTask.ID = task.ID
				tasks[i] = updatedTask

				json.NewEncoder(w).Encode(updatedTask)
				return
			}
		}

		fmt.Fprintln(w, "Task not found")
		return
	} // ← PUT ENDS HERE

	// DELETE GOES HERE
	if r.Method == http.MethodDelete {

		for i, task := range tasks {

			if id == task.ID {

				tasks = append(tasks[:i], tasks[i+1:]...)

				fmt.Fprintln(w, "Task deleted:", task.Title)
				return
			}
		}

		fmt.Fprintln(w, "Task not found")
		return
	}
}
