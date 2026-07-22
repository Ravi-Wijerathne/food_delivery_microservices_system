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

type Restaurant struct {
	ID      string  `json:"id" bson:"_id"`
	Name    string  `json:"name" bson:"name"`
	Cuisine string  `json:"cuisine" bson:"cuisine"`
	Rating  float64 `json:"rating" bson:"rating"`
	Address string  `json:"address" bson:"address"`
}

type MenuItem struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	Category    string  `json:"category" bson:"category"`
}

type Menu struct {
	RestaurantID string     `json:"restaurant_id" bson:"restaurant_id"`
	Restaurant   string     `json:"restaurant_name" bson:"restaurant_name"`
	Items        []MenuItem `json:"items" bson:"items"`
}

var restaurantClient *mongo.Client
var restaurantCollection *mongo.Collection
var menuCollection *mongo.Collection

func InitData() {
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
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	if err = c.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	restaurantClient = c
	db := restaurantClient.Database("restaurant-db")
	restaurantCollection = db.Collection("restaurants")
	menuCollection = db.Collection("menus")

	log.Println("Connected to MongoDB")

	// Seed data if empty
	count, _ := restaurantCollection.CountDocuments(ctx, bson.M{})
	if count == 0 {
		seedData()
	}
}

func seedData() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var restaurants = []interface{}{
		Restaurant{ID: "rest-1", Name: "Pizza Palace", Cuisine: "Italian", Rating: 4.5, Address: "123 Main St"},
		Restaurant{ID: "rest-2", Name: "Burger Barn", Cuisine: "American", Rating: 4.2, Address: "456 Oak Ave"},
		Restaurant{ID: "rest-3", Name: "Sushi Supreme", Cuisine: "Japanese", Rating: 4.8, Address: "789 Elm Blvd"},
		Restaurant{ID: "rest-4", Name: "Taco Town", Cuisine: "Mexican", Rating: 4.3, Address: "321 Pine Rd"},
		Restaurant{ID: "rest-5", Name: "Curry Corner", Cuisine: "Indian", Rating: 4.6, Address: "654 Maple Dr"},
	}

	var menus = []interface{}{
		Menu{
			RestaurantID: "rest-1",
			Restaurant:   "Pizza Palace",
			Items: []MenuItem{
				{ID: "item-1", Name: "Margherita Pizza", Description: "Classic tomato and mozzarella", Price: 12.99, Category: "Pizza"},
				{ID: "item-2", Name: "Pepperoni Pizza", Description: "With pepperoni slices", Price: 14.99, Category: "Pizza"},
				{ID: "item-3", Name: "Garlic Bread", Description: "Crispy with garlic butter", Price: 4.99, Category: "Sides"},
			},
		},
		Menu{
			RestaurantID: "rest-2",
			Restaurant:   "Burger Barn",
			Items: []MenuItem{
				{ID: "item-4", Name: "Classic Burger", Description: "Beef patty with lettuce and tomato", Price: 9.99, Category: "Burgers"},
				{ID: "item-5", Name: "Cheese Burger", Description: "With melted cheddar", Price: 11.99, Category: "Burgers"},
				{ID: "item-6", Name: "French Fries", Description: "Crispy golden fries", Price: 3.99, Category: "Sides"},
			},
		},
		Menu{
			RestaurantID: "rest-3",
			Restaurant:   "Sushi Supreme",
			Items: []MenuItem{
				{ID: "item-7", Name: "Salmon Roll", Description: "Fresh salmon with avocado", Price: 15.99, Category: "Rolls"},
				{ID: "item-8", Name: "Tuna Sashimi", Description: "Premium tuna slices", Price: 18.99, Category: "Sashimi"},
				{ID: "item-9", Name: "Miso Soup", Description: "Traditional miso", Price: 3.99, Category: "Soup"},
			},
		},
		Menu{
			RestaurantID: "rest-4",
			Restaurant:   "Taco Town",
			Items: []MenuItem{
				{ID: "item-10", Name: "Beef Tacos", Description: "Three tacos with seasoned beef", Price: 10.99, Category: "Tacos"},
				{ID: "item-11", Name: "Chicken Burrito", Description: "Large flour tortilla wrap", Price: 12.99, Category: "Burritos"},
				{ID: "item-12", Name: "Guacamole & Chips", Description: "Fresh guacamole", Price: 6.99, Category: "Sides"},
			},
		},
		Menu{
			RestaurantID: "rest-5",
			Restaurant:   "Curry Corner",
			Items: []MenuItem{
				{ID: "item-13", Name: "Chicken Curry", Description: "Creamy butter chicken", Price: 13.99, Category: "Curry"},
				{ID: "item-14", Name: "Palak Paneer", Description: "Spinach and cheese curry", Price: 11.99, Category: "Vegetarian"},
				{ID: "item-15", Name: "Garlic Naan", Description: "Soft bread with garlic", Price: 2.99, Category: "Breads"},
			},
		},
	}

	restaurantCollection.InsertMany(ctx, restaurants)
	menuCollection.InsertMany(ctx, menus)
	log.Println("Seeded restaurant and menu data")
}

func GetRestaurants() []Restaurant {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var restaurants []Restaurant
	cursor, err := restaurantCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching restaurants: %v", err)
		return []Restaurant{}
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &restaurants); err != nil {
		log.Printf("Error decoding restaurants: %v", err)
		return []Restaurant{}
	}

	return restaurants
}

func GetMenu(restaurantID string) (*Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var menu Menu
	err := menuCollection.FindOne(ctx, bson.M{"restaurant_id": restaurantID}).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrMenuNotFound
		}
		log.Printf("Error fetching menu: %v", err)
		return nil, err
	}

	return &menu, nil
}

var ErrMenuNotFound = &RestaurantError{"menu not found"}

type RestaurantError struct {
	Message string
}

func (e *RestaurantError) Error() string {
	return e.Message
}
