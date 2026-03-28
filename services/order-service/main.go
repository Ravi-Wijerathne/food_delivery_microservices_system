package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitStore()

	r := mux.NewRouter()
	r.HandleFunc("/orders", CreateOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{id}", GetOrderHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	log.Println("Order Service starting on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
