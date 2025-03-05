package database

import (
	"go-TodoListAPI/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	database, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
		return err
	}

	DB = database

	err = DB.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		return err
	}

	return nil
}
