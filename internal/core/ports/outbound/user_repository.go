package outbound

import (
	"context"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

// Repository is a store of user data.
type UserRepository interface {
	// GetUserByID retrieves the [User] with `id`.
	//
	// # Errors
	// 	- [NotFoundError] if no such User exists.
	GetUserByID(ctx context.Context, id int64) (*user.User, error)

	// GetUserByEmail returns a user by email.
	//
	// # Errors
	// 	- [NotFoundError] if no such User exists.
	GetUserByEmail(ctx context.Context, email user.EmailAddress) (*user.User, error)

	// CreateUser persists a new user.
	//
	// # Errors
	// 	- [ValidationError] if email or username is already taken.
	CreateUser(ctx context.Context, req *user.RegistrationRequest) (*user.User, error)
}
