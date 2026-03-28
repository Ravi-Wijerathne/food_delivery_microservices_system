package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type      string    `json:"type"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Retry     int       `json:"retry"`
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

const (
	MaxRetries = 3
	RetryDelay = 2 * time.Second
)

func NewRabbitMQ() (*RabbitMQ, error) {
	rabbitURI := os.Getenv("RABBITMQ_URI")
	if rabbitURI == "" {
		rabbitURI = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
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
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (r *RabbitMQ) PublishWithRetry(event Event) error {
	var lastErr error

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d for event %s...", attempt, event.Type)
			time.Sleep(RetryDelay)
		}

		event.Retry = attempt
		if err := r.Publish(event); err != nil {
			lastErr = err
			log.Printf("Failed to publish event (attempt %d): %v", attempt+1, err)
			continue
		}
		return nil
	}

	return fmt.Errorf("failed after %d retries: %w", MaxRetries, lastErr)
}

func (r *RabbitMQ) Publish(event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(ctx,
		"orders",
		event.Type,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("[PUBLISHED] Event: %s, OrderID: %s, Retry: %d", event.Type, event.OrderID, event.Retry)
	return nil
}

func (r *RabbitMQ) ConsumeWithRetry(eventType string, handler func([]byte) error) error {
	q, err := r.channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = r.channel.QueueBind(
		q.Name,
		eventType,
		"orders",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	log.Printf("[CONSUMER] Listening for event: %s", eventType)

	go func() {
		for d := range msgs {
			log.Printf("[RECEIVED] Event: %s, Body: %s", eventType, string(d.Body))

			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("[ERROR] Failed to unmarshal event: %v", err)
				d.Nack(false, false)
				continue
			}

			if err := handler(d.Body); err != nil {
				log.Printf("[ERROR] Handler failed: %v", err)

				if event.Retry < MaxRetries {
					log.Printf("[RETRY] Requeuing event for retry (attempt %d)", event.Retry+1)
					event.Retry++
					retryBody, _ := json.Marshal(event)
					d.Nack(false, true)

					time.Sleep(RetryDelay)
					r.channel.PublishWithContext(
						context.Background(),
						"orders",
						eventType,
						false,
						false,
						amqp.Publishing{
							ContentType:  "application/json",
							Body:         retryBody,
							DeliveryMode: amqp.Persistent,
						},
					)
				} else {
					log.Printf("[DLQ] Max retries reached, sending to dead letter")
					d.Nack(false, false)
				}
				continue
			}

			log.Printf("[SUCCESS] Event processed: %s", eventType)
			d.Ack(false)
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
