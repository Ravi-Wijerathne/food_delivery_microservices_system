package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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
	mq, err := NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.Consume("OrderCreated", handleOrderCreated)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		})
		log.Println("Payment HTTP Service starting on :8085")
		log.Fatal(http.ListenAndServe(":8085", mux))
	}()

	log.Println("Payment Service waiting for messages...")
	select {}
}

func handleOrderCreated(data []byte) error {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("Processing payment for order: %s, amount: %.2f", event.OrderID, event.Amount)

	time.Sleep(1 * time.Second)

	payment := Payment{
		ID:        fmt.Sprintf("pay-%d", time.Now().Unix()),
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		Amount:    event.Amount,
		Status:    "SUCCESS",
		CreatedAt: time.Now(),
	}
	payments[payment.ID] = payment

	log.Printf("Payment processed: %s for order: %s", payment.ID, payment.OrderID)

	event.Type = "PaymentProcessed"
	event.Timestamp = time.Now()
	mq, _ := NewRabbitMQ()
	if mq != nil {
		mq.Publish(event)
		mq.Close()
	}

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

func (r *RabbitMQ) Publish(event Event) error {
	body, _ := json.Marshal(event)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(ctx,
		"orders",
		event.Type,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published event: %s for order: %s", event.Type, event.OrderID)
	return nil
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
