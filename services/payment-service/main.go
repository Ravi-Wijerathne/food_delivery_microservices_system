package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"common"
)

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var payments = make(map[string]Payment)

func main() {
	log.Println("[PAYMENT] Starting Payment Service...")

	mq, err := common.NewRabbitMQ()
	if err != nil {
		log.Fatalf("[PAYMENT] Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.ConsumeWithRetry("OrderCreated", handleOrderCreated)
	if err != nil {
		log.Fatalf("[PAYMENT] Failed to consume: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthHandler)
		mux.HandleFunc("/payments", listPaymentsHandler)
		log.Println("[PAYMENT] HTTP Service starting on :8085")
		log.Fatal(http.ListenAndServe(":8085", mux))
	}()

	log.Println("[PAYMENT] Service ready, waiting for messages...")
	select {}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"service": "payment",
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func listPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func handleOrderCreated(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[PAYMENT] ERROR: Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf("[PAYMENT] Processing payment for order: %s, amount: $%.2f", event.OrderID, event.Amount)

	time.Sleep(500 * time.Millisecond)

	payment := Payment{
		ID:        fmt.Sprintf("pay-%d", time.Now().Unix()),
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		Amount:    event.Amount,
		Status:    "SUCCESS",
		CreatedAt: time.Now(),
	}
	payments[payment.ID] = payment

	log.Printf("[PAYMENT] SUCCESS: Payment %s processed for order: %s", payment.ID, event.OrderID)

	event.Type = "PaymentProcessed"
	event.Timestamp = time.Now()

	mq, err := common.NewRabbitMQ()
	if err != nil {
		log.Printf("[PAYMENT] ERROR: Failed to create RabbitMQ connection: %v", err)
		return err
	}
	defer mq.Close()

	if err := mq.PublishWithRetry(event); err != nil {
		log.Printf("[PAYMENT] ERROR: Failed to publish PaymentProcessed event: %v", err)
		return err
	}

	return nil
}
