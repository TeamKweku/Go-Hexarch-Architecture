package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc/middleware"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc/user"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/logger"
)

type Server struct {
	server     *grpc.Server
	port       int
	userServer *user.Server
	listener   net.Listener
	logger     logger.Logger
}

func NewServer(port int, userServer *user.Server, logger logger.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GrpcLogger(logger)),
	)
	reflection.Register(grpcServer)
	return &Server{
		server:     grpcServer,
		port:       port,
		userServer: userServer,
		logger:     logger,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = lis

	s.userServer.RegisterServer(s.server)

	s.logger.Info(
		context.Background(),
		fmt.Sprintf("gRPC server listening on :%d", s.port),
		map[string]interface{}{"port": s.port},
	)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve %w", err)
	}

	return nil
}

func (s *Server) Close() {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			s.logger.Error(context.Background(), err, "error closing listener", nil)
		}
	}
	s.server.GracefulStop()
	s.logger.Info(context.Background(), "gRPC server stopped gracefully", nil)
}
