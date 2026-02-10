package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gemini-cli/palantir/api"
	"github.com/gemini-cli/palantir/internal/server"
	"github.com/gemini-cli/palantir/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("Starting Palantir gRPC Server...")

	// Create a new BadgerStorage instance
	// TODO: Make the storage path configurable
	s, err := storage.NewBadgerStorage("C:\\Users\\saipr\\AppData\\Local\\Temp\\palantir_storage")
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer s.Close()

	// Create a gRPC server instance
	grpcServer := grpc.NewServer()

	// Register the Palantir service with the gRPC server
	palantirServer := server.NewPalantirServer(s)
	api.RegisterPalantirServer(grpcServer, palantirServer)

	// Register reflection service on gRPC server.
	// This is useful for gRPC clients to learn about the server's services.
	reflection.Register(grpcServer)

	// Define the port for the gRPC server
	grpcPort := ":50051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	// Start serving gRPC requests in a goroutine
	go func() {
		log.Printf("gRPC server listening on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Graceful shutdown
	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChan
	log.Printf("Received signal %s, shutting down gRPC server...", sig)

	// Stop the gRPC server gracefully
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")
}