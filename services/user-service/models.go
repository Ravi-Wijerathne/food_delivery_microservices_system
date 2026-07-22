package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserProfile struct {
	ID        string    `json:"id" bson:"_id"`
	Email     string    `json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	Phone     string    `json:"phone" bson:"phone"`
	Address   string    `json:"address" bson:"address"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

var client *mongo.Client
var userProfileCollection *mongo.Collection

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
		log.Fatalf("[USER] failed to connect to MongoDB: %v", err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("[USER] failed to ping MongoDB: %v", err)
	}

	client = c
	userProfileCollection = client.Database("user-db").Collection("users")
	log.Println("[USER] Connected to MongoDB")

	seedData()
}

func seedData() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := userProfileCollection.CountDocuments(ctx, bson.M{})
	if count == 0 {
		user := UserProfile{
			ID:        "user-1",
			Email:     "test@example.com",
			Name:      "Test User",
			Phone:     "1234567890",
			Address:   "123 Default Address",
			CreatedAt: time.Now(),
		}
		userProfileCollection.InsertOne(ctx, user)
		log.Println("[USER] Seeded user data")
	}
}

func GetUserByID(id string) (*UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user UserProfile
	err := userProfileCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userProfileCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"name":    user.Name,
			"phone":   user.Phone,
			"address": user.Address,
		}},
	)
	return err
}
