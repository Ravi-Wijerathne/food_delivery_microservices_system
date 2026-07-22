package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Payment struct {
	ID        string    `json:"id" bson:"_id"`
	OrderID   string    `json:"order_id" bson:"order_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

var client *mongo.Client
var paymentCollection *mongo.Collection

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("[PAYMENT] failed to connect to MongoDB: %v", err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("[PAYMENT] failed to ping MongoDB: %v", err)
	}

	client = c
	paymentCollection = client.Database("payment-db").Collection("payments")
	log.Println("[PAYMENT] Connected to MongoDB")
}

func main() {
	log.Println("[PAYMENT] Starting Payment Service...")

	InitMongoDB()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payments []Payment
	cursor, err := paymentCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &payments); err != nil {
		http.Error(w, "Failed to parse payments", http.StatusInternalServerError)
		return
	}

	if payments == nil {
		payments = []Payment{}
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := paymentCollection.InsertOne(ctx, payment)
	if err != nil {
		log.Printf("[PAYMENT] ERROR: Failed to save payment: %v", err)
		return err
	}

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
