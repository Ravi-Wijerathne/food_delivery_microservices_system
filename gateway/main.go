package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8088"
	}

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8083"
	}

	r.Use(CORSMiddleware)
	r.Use(loggingMiddleware)

	// Preflight OPTIONS for all subroutes
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.WriteHeader(http.StatusOK)
	})

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", HealthHandler).Methods("GET")
	api.Handle("/metrics", promhttp.Handler()).Methods("GET")

	api.HandleFunc("/auth/register", proxyRequest(authServiceURL)).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", proxyRequest(authServiceURL)).Methods("POST", "OPTIONS")

	api.HandleFunc("/restaurants", proxyRequest(restaurantServiceURL)).Methods("GET", "OPTIONS")
	api.HandleFunc("/menu/{id}", proxyRequest(restaurantServiceURL)).Methods("GET", "OPTIONS")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(JWTMiddleware)
	protected.HandleFunc("/users/profile", proxyRequest(userServiceURL)).Methods("GET", "PUT", "OPTIONS")
	protected.HandleFunc("/orders", proxyRequest(orderServiceURL)).Methods("POST", "OPTIONS")
	protected.HandleFunc("/orders/{id}", proxyRequest(orderServiceURL)).Methods("GET", "OPTIONS")

	// Serve static files from web directory if present
	if _, err := os.Stat("./web"); err == nil {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
	} else if _, err := os.Stat("../web"); err == nil {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("../web")))
	}

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
