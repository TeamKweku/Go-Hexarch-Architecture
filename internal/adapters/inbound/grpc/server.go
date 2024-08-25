package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/inbound/grpc/user"
)

type Server struct {
	server     *grpc.Server
	port       int
	userServer *user.Server
	listener   net.Listener
}

func NewServer(port int, userServer *user.Server) *Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	return &Server{
		server:     grpcServer,
		port:       port,
		userServer: userServer,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s.listener = lis

	s.userServer.RegisterServer(s.server)

	log.Printf("gRPC server listening on :%d", s.port)

	return s.server.Serve(lis)
}

func (s *Server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
	s.server.GracefulStop()
	log.Println("gRPC server stopped gracefully")
}
