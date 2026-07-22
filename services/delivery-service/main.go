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

type Delivery struct {
	ID         string    `json:"id" bson:"_id"`
	OrderID    string    `json:"order_id" bson:"order_id"`
	AgentName  string    `json:"agent_name" bson:"agent_name"`
	Status     string    `json:"status" bson:"status"`
	AssignedAt time.Time `json:"assigned_at" bson:"assigned_at"`
}

var client *mongo.Client
var deliveryCollection *mongo.Collection

var agentNames = []string{"John", "Mike", "Sarah", "David", "Emma"}

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
		log.Fatalf("[DELIVERY] failed to connect to MongoDB: %v", err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("[DELIVERY] failed to ping MongoDB: %v", err)
	}

	client = c
	deliveryCollection = client.Database("delivery-db").Collection("deliveries")
	log.Println("[DELIVERY] Connected to MongoDB")
}

func main() {
	log.Println("[DELIVERY] Starting Delivery Service...")

	InitMongoDB()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var deliveries []Delivery
	cursor, err := deliveryCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch deliveries", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &deliveries); err != nil {
		http.Error(w, "Failed to parse deliveries", http.StatusInternalServerError)
		return
	}

	if deliveries == nil {
		deliveries = []Delivery{}
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := deliveryCollection.InsertOne(ctx, delivery)
	if err != nil {
		log.Printf("[DELIVERY] ERROR: Failed to save delivery: %v", err)
		return err
	}

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
