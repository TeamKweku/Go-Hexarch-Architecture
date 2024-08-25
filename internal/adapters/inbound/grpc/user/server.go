package user

import (
	pb "github.com/TeamKweku/code-odessey-hex-arch-proto/protogen/go/user"
	"google.golang.org/grpc"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userService inbound.UserService
}

func NewServer(userService inbound.UserService) *Server {
	return &Server{userService: userService}
}

func (s *Server) RegisterServer(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
}
