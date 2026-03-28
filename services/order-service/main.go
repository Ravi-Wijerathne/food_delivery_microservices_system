package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/food_delivery_microservices_system/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	InitStore()

	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/orders", CreateOrderHandler).Methods("POST")
		r.HandleFunc("/orders/{id}", GetOrderHandler).Methods("GET")
		r.HandleFunc("/health", HealthHandler).Methods("GET")

		log.Println("Order HTTP Service starting on :8083")
		log.Fatal(http.ListenAndServe(":8083", r))
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8084")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterOrderServiceServer(grpcServer, &orderGrpcServer{})

		log.Println("Order gRPC Service starting on :8084")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	<-sigCh
	log.Println("Shutting down...")
}

type orderGrpcServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *orderGrpcServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var total float64
	for _, item := range req.Items {
		total += item.Price * float64(item.Quantity)
	}

	order := Order{
		ID:              generateOrderID(),
		UserID:          req.UserId,
		RestaurantID:    req.RestaurantId,
		Items:           convertItems(req.Items),
		TotalAmount:     total,
		Status:          OrderStatusCreated,
		DeliveryAddress: req.DeliveryAddress,
	}

	if err := SaveOrder(order); err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:              order.ID,
		UserId:          order.UserID,
		RestaurantId:    order.RestaurantID,
		Items:           convertItemsToProto(order.Items),
		TotalAmount:     order.TotalAmount,
		Status:          string(order.Status),
		DeliveryAddress: order.DeliveryAddress,
		CreatedAt:       order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *orderGrpcServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := GetOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetOrderResponse{
		Id:              order.ID,
		UserId:          order.UserID,
		RestaurantId:    order.RestaurantID,
		Items:           convertItemsToProto(order.Items),
		TotalAmount:     order.TotalAmount,
		Status:          string(order.Status),
		DeliveryAddress: order.DeliveryAddress,
		CreatedAt:       order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *orderGrpcServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	err := UpdateOrderStatus(req.OrderId, OrderStatus(req.Status))
	if err != nil {
		return &pb.UpdateOrderStatusResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.UpdateOrderStatusResponse{Success: true, Message: "Status updated"}, nil
}

func convertItems(items []*pb.OrderItem) []OrderItem {
	result := make([]OrderItem, len(items))
	for i, item := range items {
		result[i] = OrderItem{
			MenuItemID: item.MenuItemId,
			Name:       item.Name,
			Quantity:   int(item.Quantity),
			Price:      item.Price,
		}
	}
	return result
}

func convertItemsToProto(items []OrderItem) []*pb.OrderItem {
	result := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		result[i] = &pb.OrderItem{
			MenuItemId: item.MenuItemID,
			Name:       item.Name,
			Quantity:   int32(item.Quantity),
			Price:      item.Price,
		}
	}
	return result
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
