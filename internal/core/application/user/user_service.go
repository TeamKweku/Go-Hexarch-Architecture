package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	domainUser "github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
)

// service that implements and satisfies the inbound service interface
type UserService struct {
	repo               outbound.UserRepository
	passwordComparator domainUser.PasswordComparator
}

func NewUserService(repo outbound.UserRepository) *UserService {
	return &UserService{
		repo:               repo,
		passwordComparator: domainUser.BcryptCompare,
	}
}

func (us *UserService) Register(
	ctx context.Context,
	req *domainUser.RegistrationRequest,
) (*domainUser.User, error) {
	user, err := us.repo.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user from %#v: %w", req, err)
	}

	return user, nil
}

func (us *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*domainUser.User, error) {
	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user with ID %s: %w", id, err)
	}
	return user, nil
}

func (us *UserService) Authenticate(
	ctx context.Context,
	req *domainUser.LoginRequest,
) (*domainUser.User, error) {
	user, err := us.repo.GetUserByEmail(ctx, req.Email())
	if err != nil {
		var notFoundErr *domainUser.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, &domainUser.AuthError{Cause: err}
		}
		return nil, err
	}
	if err := us.passwordComparator(user.PasswordHash(), req.PasswordCandidate()); err != nil {
		return nil, &domainUser.AuthError{Cause: err}
	}

	return user, nil
}

func (us *UserService) UpdateUser(
	ctx context.Context,
	req *domainUser.UpdateRequest,
) (*domainUser.User, error) {
	user, err := us.repo.UpdateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update user with ID %s: %w", req.UserID().String(), err)
	}

	return user, nil
}
