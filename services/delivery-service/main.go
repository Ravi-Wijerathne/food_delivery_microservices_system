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
	mq, err := NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	err = mq.Consume("PaymentProcessed", handlePaymentProcessed)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		})
		log.Println("Delivery HTTP Service starting on :8086")
		log.Fatal(http.ListenAndServe(":8086", mux))
	}()

	log.Println("Delivery Service waiting for messages...")
	select {}
}

func handlePaymentProcessed(data []byte) error {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("Assigning delivery for order: %s", event.OrderID)

	time.Sleep(1 * time.Second)

	agent := agentNames[time.Now().Unix()%int64(len(agentNames))]

	delivery := Delivery{
		ID:         fmt.Sprintf("del-%d", time.Now().Unix()),
		OrderID:    event.OrderID,
		AgentName:  agent,
		Status:     "ASSIGNED",
		AssignedAt: time.Now(),
	}
	deliveries[delivery.ID] = delivery

	log.Printf("Delivery assigned: %s to %s for order: %s", delivery.ID, agent, event.OrderID)

	event.Type = "DeliveryAssigned"
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
