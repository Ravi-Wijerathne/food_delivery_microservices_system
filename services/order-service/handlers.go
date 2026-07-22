package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"common"
	pb "github.com/food_delivery_microservices_system/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// 1. Verify User via gRPC
	userSvcURL := os.Getenv("USER_SERVICE_GRPC_URL")
	if userSvcURL == "" {
		userSvcURL = "localhost:9088" // default local
	}
	userConn, err := grpc.Dial(userSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to User Service: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer userConn.Close()

	userClient := pb.NewUserServiceClient(userConn)
	_, err = userClient.GetUser(context.Background(), &pb.GetUserRequest{Id: req.UserID})
	if err != nil {
		log.Printf("User not found or error: %v", err)
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// 2. Verify Restaurant Menu Items via gRPC and calculate total
	restSvcURL := os.Getenv("RESTAURANT_SERVICE_GRPC_URL")
	if restSvcURL == "" {
		restSvcURL = "localhost:9082" // default local
	}
	restConn, err := grpc.Dial(restSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to Restaurant Service: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer restConn.Close()

	restClient := pb.NewRestaurantServiceClient(restConn)
	var total float64
	for i, item := range req.Items {
		menuItem, err := restClient.GetMenuItem(context.Background(), &pb.GetMenuItemRequest{
			RestaurantId: req.RestaurantID,
			ItemId:       item.MenuItemID,
		})
		if err != nil {
			log.Printf("Menu item not found: %v", err)
			http.Error(w, "Invalid Menu Item ID", http.StatusBadRequest)
			return
		}
		// Trust the price from the server, not the client
		req.Items[i].Price = menuItem.Price
		req.Items[i].Name = menuItem.Name
		total += menuItem.Price * float64(item.Quantity)
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
