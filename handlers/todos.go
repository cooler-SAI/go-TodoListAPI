package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go-TodoListAPI/database"
	"go-TodoListAPI/internal/auth"
	"go-TodoListAPI/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.UserID = userID
	if err := database.DB.Create(&todo).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err2 := json.NewEncoder(w).Encode(todo)
	if err2 != nil {
		return
	}
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var todos []models.Todo
	if err := database.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		return
	}
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err2 := json.NewEncoder(w).Encode(todo)
	if err2 != nil {
		return
	}
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	var updatedTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.Title = updatedTodo.Title
	todo.Description = updatedTodo.Description
	if err := database.DB.Save(&todo).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err2 := json.NewEncoder(w).Encode(todo)
	if err2 != nil {
		return
	}
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Todo{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // Удаляем "Bearer " из заголовка

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
