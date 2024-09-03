package user

import (
	"context"
	"errors"

	pb "github.com/TeamKweku/code-odessey-hex-arch-proto/protogen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	domainUser "github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userService inbound.UserService
	config      config.Config
	authService inbound.TokenService
}

func NewServer(
	userService inbound.UserService,
	cfg config.Config,
	authService inbound.TokenService,
) *Server {
	return &Server{
		userService: userService,
		authService: authService,
		config:      cfg,
	}
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

// Authenticate user function to login user
func (s *Server) Authenticate(
	ctx context.Context,
	req *pb.LoginUserRequest,
) (*pb.LoginUserResponse, error) {
	loginReq, err := domainUser.ParseLoginRequest(req.GetEmail(), req.GetPassword())
	if err != nil {
		var validatorErrors domainUser.ValidationErrors
		if errors.As(err, &validatorErrors) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", validatorErrors)
		}
		return nil, status.Errorf(codes.Internal, "failed to parse login request: %v", err)
	}

	// authenticate the user
	authenticatedUser, err := s.userService.Authenticate(ctx, loginReq)
	if err != nil {
		var authErr *domainUser.AuthError
		if errors.As(err, &authErr) {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to authenticate user: %v", err)
	}

	// generate access token
	accessToken, _, err := s.authService.CreateToken(authenticatedUser, s.config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v", err)
	}

	return &pb.LoginUserResponse{
		User: &pb.User{
			Id:                authenticatedUser.ID().String(),
			Username:          authenticatedUser.Username().String(),
			Email:             authenticatedUser.Email().String(),
			Role:              authenticatedUser.Role().String(),
			PasswordChangedAt: timestamppb.New(authenticatedUser.PasswordChangedAt()),
			CreatedAt:         timestamppb.New(authenticatedUser.CreatedAt()),
			UpdatedAt:         timestamppb.New(authenticatedUser.UpdatedAt()),
		},
		AccessToken: accessToken,
	}, nil
}
