package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/food_delivery_microservices_system/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type restaurantGrpcServer struct {
	pb.UnimplementedRestaurantServiceServer
}

func (s *restaurantGrpcServer) GetMenuItem(ctx context.Context, req *pb.GetMenuItemRequest) (*pb.GetMenuItemResponse, error) {
	slog.Info("gRPC GetMenuItem called", "restaurant_id", req.RestaurantId, "item_id", req.ItemId)
	
	menu, err := GetMenu(req.RestaurantId)
	if err != nil {
		return nil, err
	}

	for _, item := range menu.Items {
		if item.ID == req.ItemId {
			return &pb.GetMenuItemResponse{
				Id:    item.ID,
				Name:  item.Name,
				Price: item.Price,
			}, nil
		}
	}

	return nil, ErrMenuNotFound
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Restaurant Service")
	InitData()

	// Start HTTP Server
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/restaurants", GetRestaurantsHandler).Methods("GET")
		r.HandleFunc("/menu/{id}", GetMenuHandler).Methods("GET")
		r.HandleFunc("/health", HealthHandler).Methods("GET")
		
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "8082"
		}
		
		slog.Info("Restaurant HTTP Service starting", "port", httpPort)
		if err := http.ListenAndServe(":"+httpPort, r); err != nil {
			slog.Error("HTTP Server failed", "error", err)
		}
	}()

	// Start gRPC Server
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "9082"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			os.Exit(1)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterRestaurantServiceServer(grpcServer, &restaurantGrpcServer{})

		slog.Info("Restaurant gRPC Service starting", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC Server failed", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	<-sigCh

	slog.Info("Shutting down Restaurant Service")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"service": "restaurant",
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}
