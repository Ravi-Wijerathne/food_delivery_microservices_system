package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"common"
	"github.com/gorilla/mux"
)

type CreateOrderRequest struct {
	UserID          string      `json:"user_id"`
	RestaurantID    string      `json:"restaurant_id"`
	Items           []OrderItem `json:"items"`
	DeliveryAddress string      `json:"delivery_address"`
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.RestaurantID == "" || len(req.Items) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var total float64
	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}

	order := Order{
		ID:              generateOrderID(),
		UserID:          req.UserID,
		RestaurantID:    req.RestaurantID,
		Items:           req.Items,
		TotalAmount:     total,
		Status:          OrderStatusCreated,
		DeliveryAddress: req.DeliveryAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := SaveOrder(order); err != nil {
		log.Printf("Error saving order: %v", err)
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	log.Printf("[ORDER] Order created: %s", order.ID)

	if mq != nil {
		event := common.Event{
			Type:      "OrderCreated",
			OrderID:   order.ID,
			UserID:    order.UserID,
			Amount:    order.TotalAmount,
			Status:    string(order.Status),
			Timestamp: time.Now(),
		}
		if err := mq.PublishWithRetry(event); err != nil {
			log.Printf("[ORDER] Failed to publish event: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, err := GetOrder(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
