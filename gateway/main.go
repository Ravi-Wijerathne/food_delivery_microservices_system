package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8081"
	}

	restaurantServiceURL := os.Getenv("RESTAURANT_SERVICE_URL")
	if restaurantServiceURL == "" {
		restaurantServiceURL = "http://localhost:8082"
	}

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8083"
	}

	r.Use(loggingMiddleware)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", HealthHandler).Methods("GET")

	api.HandleFunc("/auth/register", proxyRequest(authServiceURL)).Methods("POST")
	api.HandleFunc("/auth/login", proxyRequest(authServiceURL)).Methods("POST")

	api.HandleFunc("/restaurants", proxyRequest(restaurantServiceURL)).Methods("GET")
	api.HandleFunc("/menu/{id}", proxyRequest(restaurantServiceURL)).Methods("GET")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(JWTMiddleware)
	protected.HandleFunc("/orders", proxyRequest(orderServiceURL)).Methods("POST")
	protected.HandleFunc("/orders/{id}", proxyRequest(orderServiceURL)).Methods("GET")

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
