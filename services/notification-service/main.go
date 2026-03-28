package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	mq, err := NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.Consume("OrderCreated", handleOrderCreated)
	if err != nil {
		log.Fatalf("Failed to consume OrderCreated: %v", err)
	}

	err = mq.Consume("PaymentProcessed", handlePaymentProcessed)
	if err != nil {
		log.Fatalf("Failed to consume PaymentProcessed: %v", err)
	}

	err = mq.Consume("DeliveryAssigned", handleDeliveryAssigned)
	if err != nil {
		log.Fatalf("Failed to consume DeliveryAssigned: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		})
		log.Println("Notification HTTP Service starting on :8087")
		log.Fatal(http.ListenAndServe(":8087", mux))
	}()

	log.Println("Notification Service waiting for messages...")
	select {}
}

func handleOrderCreated(data []byte) error {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Order Created: Order %s has been placed. Amount: $%.2f\n", event.OrderID, event.Amount)
	return nil
}

func handlePaymentProcessed(data []byte) error {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Payment Processed: Order %s payment of $%.2f is complete\n", event.OrderID, event.Amount)
	return nil
}

func handleDeliveryAssigned(data []byte) error {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	fmt.Printf("[NOTIFICATION] Delivery Assigned: A driver has been assigned to order %s\n", event.OrderID)
	return nil
}

type Event struct {
	Type      string    `json:"type"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"orders",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (r *RabbitMQ) Consume(eventType string, handler func([]byte) error) error {
	q, err := r.channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		q.Name,
		eventType,
		"orders",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
