package inbound

import (
	"context"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

type UserService interface {
	// Register a new user.
	//
	// # Errors
	// 	- [ValidationError] if email or username is already taken.
	Register(ctx context.Context, req *user.RegistrationRequest) (*user.User, error)

	// GetUser a user by ID.
	//
	// # Errors
	// 	- [NotFoundError] if no such User exists.
	GetUser(ctx context.Context, id int64) (*user.User, error)

	// Authenticate a user, returning the authenticated [User] if successful.
	//
	// # Errors
	//	- [AuthError].
	Authenticate(ctx context.Context, req *user.LoginRequest) (*user.User, error)
}
