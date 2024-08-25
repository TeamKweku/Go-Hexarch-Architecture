package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc"
	userGRPC "github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/postgres"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/application/user"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// initialize database connection
	postgresURL := postgres.NewURL(cfg)
	postgresAdapter, err := postgres.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		if err := postgresAdapter.Close(); err != nil {
			log.Printf("error closing postgres client %v", err)
		}
	}()

	userService := user.NewUserService(postgresAdapter)

	userServer := userGRPC.NewServer(userService)
	grpcServer := grpc.NewServer(8080, userServer)

	// start gRPC server
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	grpcServer.Close()
	log.Println("server exited properly")
}
