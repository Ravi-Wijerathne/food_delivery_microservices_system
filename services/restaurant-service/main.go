package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitData()

	r := mux.NewRouter()
	r.HandleFunc("/restaurants", GetRestaurantsHandler).Methods("GET")
	r.HandleFunc("/menu/{id}", GetMenuHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	log.Println("Restaurant Service starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
