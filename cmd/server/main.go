package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc"
	userGRPC "github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/postgres"
	authService "github.com/teamkweku/code-odessey-hex-arch/internal/core/application/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/application/session"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/application/user"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/logger"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	isPrettyPrint := cfg.Environment == "development"

	// initialize logger
	zeroLogger := logger.NewZerologLogger(isPrettyPrint)

	// initialize database connection
	postgresURL := postgres.NewURL(cfg)
	postgresAdapter, err := postgres.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		if err = postgresAdapter.Close(); err != nil {
			log.Printf("error closing postgres client %v", err)
		}
	}()

	// create sessions service
	sessionSrv := session.NewSessionService(postgresAdapter)
	userService := user.NewUserService(postgresAdapter)

	// initialize token service
	tokenService, err := auth.NewPasetoToken()
	if err != nil {
		log.Fatalf("failed to create token service: %v", err)
	}

	// initialize the auth service
	authService := authService.NewAuthService(tokenService, sessionSrv)

	userServer := userGRPC.NewServer(userService, cfg, authService, sessionSrv)

	// Convert RPCPort to int
	port, err := strconv.Atoi(cfg.RPCPort)
	if err != nil {
		log.Fatalf("invalid port number: %v", err)
	}

	grpcServer := grpc.NewServer(port, userServer, zeroLogger)

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
