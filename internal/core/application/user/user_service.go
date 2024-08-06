package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
)

// service that implements and satisfies the inbound service interface
type UserService struct {
	repo               outbound.UserRepository
	passwordComparator user.PasswordComparator
}

func NewUserService(repo outbound.UserRepository) *UserService {
	return &UserService{
		repo:               repo,
		passwordComparator: user.BcryptCompare,
	}
}

func (us *UserService) Register(ctx context.Context, req user.RegistrationRequest) (*user.User, error) {
	// TODO
	return nil, nil
}

func (us *UserService) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	// TODO
	return nil, nil
}

func (us *UserService) Authenticate(ctx context.Context, req *user.LoginRequest) (*user.User, error) {
	// TODO
	return nil, nil
}

func (us *UserService) UpdateUser(ctx context.Context, req *user.UpdateRequest) (*user.User, error) {
	// TODO
	return nil, nil
}
