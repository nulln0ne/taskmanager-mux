package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func CreateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		json.NewDecoder(r.Body).Decode(&task)

		query := "INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id"
		err := db.QueryRow(query, task.Title, task.Description, task.Completed).Scan(&task.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(task)
	}
}

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, description, completed FROM tasks")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
			tasks = append(tasks, task)
		}

		json.NewEncoder(w).Encode(tasks)
	}
}

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var task Task
		query := "SELECT id, title, description, completed FROM tasks WHERE id=$1"
		err := db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(task)
	}
}

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		var task Task
		json.NewDecoder(r.Body).Decode(&task)

		query := "UPDATE tasks SET title=$1, description=$2, completed=$3 WHERE id=$4"
		_, err := db.Exec(query, task.Title, task.Description, task.Completed, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		task.ID = id
		json.NewEncoder(w).Encode(task)
	}
}

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		query := "DELETE FROM tasks WHERE id=$1"
		_, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
