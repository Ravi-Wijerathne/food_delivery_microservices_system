package main

import (
	"errors"
	"log"
	"time"
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
	ID              string      `json:"id"`
	UserID          string      `json:"user_id"`
	RestaurantID    string      `json:"restaurant_id"`
	Items           []OrderItem `json:"items"`
	TotalAmount     float64     `json:"total_amount"`
	Status          OrderStatus `json:"status"`
	DeliveryAddress string      `json:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	MenuItemID string  `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type OrderStore struct {
	orders map[string]Order
}

var orderStore *OrderStore

func InitStore() {
	orderStore = &OrderStore{
		orders: make(map[string]Order),
	}
}

func SaveOrder(order Order) error {
	orderStore.orders[order.ID] = order
	log.Printf("Order created: %s", order.ID)
	return nil
}

func GetOrder(orderID string) (Order, error) {
	if order, exists := orderStore.orders[orderID]; exists {
		return order, nil
	}
	return Order{}, ErrOrderNotFound
}

func UpdateOrderStatus(orderID string, status OrderStatus) error {
	if order, exists := orderStore.orders[orderID]; exists {
		order.Status = status
		order.UpdatedAt = time.Now()
		orderStore.orders[orderID] = order
		return nil
	}
	return ErrOrderNotFound
}

var ErrOrderNotFound = errors.New("order not found")

func generateOrderID() string {
	return "order-" + time.Now().Format("20060102150405")
}
