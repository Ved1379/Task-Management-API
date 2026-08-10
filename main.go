package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func Homehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Task Management API")
}

func main() {
	fmt.Println("Server started on port 8080")
	http.HandleFunc("/", Homehandler)
	http.HandleFunc("/tasks", taskHandler)
	http.ListenAndServe(":8080", nil)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "No tasks yet")
	}

	if r.Method == http.MethodPost {

		var task Task

		err := json.NewDecoder(r.Body).Decode(&task)

		if err != nil {
			fmt.Fprintln(w, "Invalid JSON")
			return
		}

		fmt.Fprintln(w, "Task created:", task.Title)
	}
}
