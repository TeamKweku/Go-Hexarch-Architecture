package outbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

// Repository is a store of user data.
type UserRepository interface {
	// GetUserByID retrieves the [User] with `id`.
	//
	// # Errors
	// 	- [NotFoundError] if no such User exists.
	GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error)

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
	// UpdateUser updates an existing user.
	//
	// # Errors
	// 	- [NotFoundError] if no such User exists.
	// 	- [ValidationError] if email is already taken.
	//  - [ConcurrentModificationError] if the user has been modified since the last
	//    read.
	UpdateUser(ctx context.Context, req *user.UpdateRequest) (*user.User, error)
}
