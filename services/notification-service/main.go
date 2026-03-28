package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"common"
)

func main() {
	log.Println("[NOTIFICATION] Starting Notification Service...")

	mq, err := common.NewRabbitMQ()
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.ConsumeWithRetry("OrderCreated", handleOrderCreated)
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to consume OrderCreated: %v", err)
	}

	err = mq.ConsumeWithRetry("PaymentProcessed", handlePaymentProcessed)
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to consume PaymentProcessed: %v", err)
	}

	err = mq.ConsumeWithRetry("DeliveryAssigned", handleDeliveryAssigned)
	if err != nil {
		log.Fatalf("[NOTIFICATION] Failed to consume DeliveryAssigned: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthHandler)
		log.Println("[NOTIFICATION] HTTP Service starting on :8087")
		log.Fatal(http.ListenAndServe(":8087", mux))
	}()

	log.Println("[NOTIFICATION] Service ready, waiting for messages...")
	select {}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"service": "notification",
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func handleOrderCreated(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Order Created: Order %s has been placed. Amount: $%.2f\n", event.OrderID, event.Amount)
	return nil
}

func handlePaymentProcessed(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Payment Processed: Order %s payment of $%.2f is complete\n", event.OrderID, event.Amount)
	return nil
}

func handleDeliveryAssigned(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Delivery Assigned: A driver has been assigned to order %s\n", event.OrderID)
	return nil
}
