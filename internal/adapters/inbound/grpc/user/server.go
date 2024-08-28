package user

import (
	"context"
	"errors"

	pb "github.com/TeamKweku/code-odessey-hex-arch-proto/protogen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	domainUser "github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
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

func (s *Server) Register(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	domainReq, err := domainUser.ParseRegistrationRequest(
		req.GetUsername(),
		req.GetEmail(),
		req.GetPasswordHash(),
	)
	if err != nil {
		var validationErrors domainUser.ValidationErrors
		if errors.As(err, &validationErrors) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", validationErrors)
		}
		return nil, status.Errorf(codes.Internal, "failed to parse registration request: %v", err)
	}

	registeredUser, err := s.userService.Register(context.Background(), domainReq)
	if err != nil {
		// TODO
		var validatorErr *domainUser.ValidationError
		if errors.As(err, &validatorErr) {
			return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", validatorErr)
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterUserResponse{
		User: &pb.User{
			Id:                registeredUser.ID().String(),
			Username:          registeredUser.Username().String(),
			Email:             registeredUser.Email().String(),
			Role:              registeredUser.Role().String(),
			PasswordChangedAt: timestamppb.New(registeredUser.PasswordChangedAt()),
			CreatedAt:         timestamppb.New(registeredUser.CreatedAt()),
			UpdatedAt:         timestamppb.New(registeredUser.UpdatedAt()),
		},
	}, nil
}
