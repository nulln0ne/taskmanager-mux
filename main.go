package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"taskmanager/handlers"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v\n", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/tasks", handlers.CreateTask(db)).Methods("POST")
	r.HandleFunc("/tasks", handlers.GetTasks(db)).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.GetTask(db)).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.UpdateTask(db)).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTask(db)).Methods("DELETE")

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
