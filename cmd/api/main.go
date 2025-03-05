package main

import (
	"fmt"
	"go-TodoListAPI/database"
	"log"
)

func main() {

	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database")
}
