package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"common"
)

type Delivery struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	AgentName  string    `json:"agent_name"`
	Status     string    `json:"status"`
	AssignedAt time.Time `json:"assigned_at"`
}

var deliveries = make(map[string]Delivery)

var agentNames = []string{"John", "Mike", "Sarah", "David", "Emma"}

func main() {
	log.Println("[DELIVERY] Starting Delivery Service...")

	mq, err := common.NewRabbitMQ()
	if err != nil {
		log.Fatalf("[DELIVERY] Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.ConsumeWithRetry("PaymentProcessed", handlePaymentProcessed)
	if err != nil {
		log.Fatalf("[DELIVERY] Failed to consume: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthHandler)
		mux.HandleFunc("/deliveries", listDeliveriesHandler)
		log.Println("[DELIVERY] HTTP Service starting on :8086")
		log.Fatal(http.ListenAndServe(":8086", mux))
	}()

	log.Println("[DELIVERY] Service ready, waiting for messages...")
	select {}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"service": "delivery",
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func listDeliveriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

func handlePaymentProcessed(data []byte) error {
	var event common.Event
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[DELIVERY] ERROR: Failed to unmarshal event: %v", err)
		return err
	}

	log.Printf("[DELIVERY] Assigning delivery for order: %s", event.OrderID)

	time.Sleep(500 * time.Millisecond)

	agent := agentNames[time.Now().Unix()%int64(len(agentNames))]

	delivery := Delivery{
		ID:         fmt.Sprintf("del-%d", time.Now().Unix()),
		OrderID:    event.OrderID,
		AgentName:  agent,
		Status:     "ASSIGNED",
		AssignedAt: time.Now(),
	}
	deliveries[delivery.ID] = delivery

	log.Printf("[DELIVERY] SUCCESS: Delivery %s assigned to %s for order: %s", delivery.ID, agent, event.OrderID)

	event.Type = "DeliveryAssigned"
	event.Timestamp = time.Now()

	mq, err := common.NewRabbitMQ()
	if err != nil {
		log.Printf("[DELIVERY] ERROR: Failed to create RabbitMQ connection: %v", err)
		return err
	}
	defer mq.Close()

	if err := mq.PublishWithRetry(event); err != nil {
		log.Printf("[DELIVERY] ERROR: Failed to publish DeliveryAssigned event: %v", err)
		return err
	}

	return nil
}
