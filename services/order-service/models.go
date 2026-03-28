package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID              string      `json:"id" bson:"_id"`
	UserID          string      `json:"user_id" bson:"user_id"`
	RestaurantID    string      `json:"restaurant_id" bson:"restaurant_id"`
	Items           []OrderItem `json:"items" bson:"items"`
	TotalAmount     float64     `json:"total_amount" bson:"total_amount"`
	Status          OrderStatus `json:"status" bson:"status"`
	DeliveryAddress string      `json:"delivery_address" bson:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" bson:"updated_at"`
}

type OrderItem struct {
	MenuItemID string  `json:"menu_item_id" bson:"menu_item_id"`
	Name       string  `json:"name" bson:"name"`
	Quantity   int     `json:"quantity" bson:"quantity"`
	Price      float64 `json:"price" bson:"price"`
}

var orderClient *mongo.Client
var orderCollection *mongo.Collection

func InitMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb://localhost:27017").SetServerAPIOptions(serverAPI)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	orderClient = c
	orderCollection = orderClient.Database("order-db").Collection("orders")

	fmt.Println("Connected to MongoDB")
	return nil
}

func SaveOrder(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := orderCollection.InsertOne(ctx, order)
	return err
}

func GetOrder(orderID string) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order Order
	err := orderCollection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	return order, err
}

func UpdateOrderStatus(orderID string, status OrderStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := orderCollection.UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}

var ErrOrderNotFound = errors.New("order not found")

func generateOrderID() string {
	return "order-" + time.Now().Format("20060102150405")
}
