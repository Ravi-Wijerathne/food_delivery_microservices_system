package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", HealthHandler).Methods("GET")

	api.HandleFunc("/auth/register", proxyRequest("http://localhost:8081")).Methods("POST")
	api.HandleFunc("/auth/login", proxyRequest("http://localhost:8081")).Methods("POST")

	api.HandleFunc("/restaurants", proxyRequest("http://localhost:8082")).Methods("GET")
	api.HandleFunc("/menu/{id}", proxyRequest("http://localhost:8082")).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(JWTMiddleware)
	protected.HandleFunc("/orders", proxyRequest("http://localhost:8083")).Methods("POST")
	protected.HandleFunc("/orders/{id}", proxyRequest("http://localhost:8083")).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
