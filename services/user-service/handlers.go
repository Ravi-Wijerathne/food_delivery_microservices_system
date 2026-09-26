package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log/slog"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/profile", GetProfileHandler).Methods("GET")
	r.HandleFunc("/profile", UpdateProfileHandler).Methods("PUT")
	r.HandleFunc("/users/profile", GetProfileHandler).Methods("GET")
	r.HandleFunc("/users/profile", UpdateProfileHandler).Methods("PUT")
	r.HandleFunc("/health", HealthHandler).Methods("GET")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"service": "user",
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Mock fetching from auth token
	userID := "user-1" 

	user, err := GetUserByID(userID)
	if err != nil {
		slog.Error("Failed to fetch user profile", "error", err, "user_id", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For simple implementation, force user-1
	req.ID = "user-1"

	if err := UpdateUser(req); err != nil {
		slog.Error("Failed to update user", "error", err, "user_id", req.ID)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
