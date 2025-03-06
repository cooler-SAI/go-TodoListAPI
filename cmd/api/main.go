package main

import (
	"fmt"
	"go-TodoListAPI/database"
	"go-TodoListAPI/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	r.HandleFunc("/todos", handlers.AuthMiddleware(handlers.CreateTodo)).Methods("POST")
	r.HandleFunc("/todos", handlers.AuthMiddleware(handlers.GetTodos)).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.AuthMiddleware(handlers.GetTodo)).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.AuthMiddleware(handlers.UpdateTodo)).Methods("PUT")
	r.HandleFunc("/todos/{id}", handlers.AuthMiddleware(handlers.DeleteTodo)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
