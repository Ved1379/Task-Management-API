package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var tasks []Task
var nextID = 1
var db *sql.DB

func Homehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Task Management API")
}

func main() {
	err := godotenv.Overload()

	if err != nil {
		fmt.Println("Error loading .env")
		return
	}

	connectionString := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	fmt.Println("HOST:", os.Getenv("DB_HOST"))
	fmt.Println("PORT:", os.Getenv("DB_PORT"))
	fmt.Println("USER:", os.Getenv("DB_USER"))
	fmt.Println("DB NAME:", os.Getenv("DB_NAME"))
	fmt.Println("SSL MODE: disable")

	db, err = sql.Open("postgres", connectionString)
	db, err = sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Println("Database connection error", err)
		return
	}

	err = db.Ping()

	if err != nil {
		fmt.Println("Database is not reaachable", err)
		return
	}
	fmt.Println("Database connected Successfully")
	fmt.Println("Server started on port 8080")
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/tasks/", taskHandler)
	http.ListenAndServe(":8080", nil)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if r.URL.Path == "/tasks" {

		if r.Method == http.MethodGet {

			rows, err := db.Query("SELECT id, title, description FROM tasks")

			if err != nil {
				http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var tasks []Task

			for rows.Next() {

				var task Task

				err := rows.Scan(
					&task.ID,
					&task.Title,
					&task.Description,
				)
				if err != nil {
					http.Error(w, "Error in reading tasks", http.StatusInternalServerError)
					return
				}
				tasks = append(tasks, task)
			}

			if err := rows.Err(); err != nil {
				http.Error(w, "Error reading tasks", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)

			return

		}

		if r.Method == http.MethodPost {

			var task Task

			err := json.NewDecoder(r.Body).Decode(&task)

			if err != nil {
				fmt.Fprintln(w, "Invalid JSON")
				return
			}
			query := "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id"

			fmt.Println("SQL:", query)
			err = db.QueryRow(
				query,
				task.Title,
				task.Description,
			).Scan(&task.ID)

			if err != nil {
				fmt.Fprintln(w, "Error creating task", err)
				return
			}
			json.NewEncoder(w).Encode(task)
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
			fmt.Fprintln(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		query := `
		UPDATE tasks
		SET title = $1, description = $2
		WHERE id = $3
		RETURNING id, title, description
		`
		err = db.QueryRow(
			query,
			updatedTask.Title,
			updatedTask.Description,
			id,
		).Scan(
			&updatedTask.ID,
			&updatedTask.Title,
			&updatedTask.Description,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Error updating task", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(updatedTask)

		return
	}

	if r.Method == http.MethodDelete {

		query := "DELETE FROM tasks WHERE id = $1"

		result, err := db.Exec(query, id)

		if err != nil {
			http.Error(w, "Error deleting task", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			http.Error(w, "Error checking deleted task", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Task deleted successfully",
		})

		return
	}
}
