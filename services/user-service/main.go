package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/food_delivery_microservices_system/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type userGrpcServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *userGrpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	slog.Info("gRPC GetUser called", "user_id", req.Id)
	user, err := GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting User Service")

	InitMongoDB()

	// Start HTTP Server
	go func() {
		r := mux.NewRouter()
		RegisterRoutes(r)
		
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "8088"
		}
		
		slog.Info("User HTTP Service starting", "port", httpPort)
		if err := http.ListenAndServe(":"+httpPort, r); err != nil {
			slog.Error("HTTP Server failed", "error", err)
		}
	}()

	// Start gRPC Server
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "9088"
		}

		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			slog.Error("Failed to listen for gRPC", "error", err)
			os.Exit(1)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(grpcServer, &userGrpcServer{})

		slog.Info("User gRPC Service starting", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC Server failed", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT)
	<-sigCh

	slog.Info("Shutting down User Service")
}
