package outbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *auth.Sessions) error
	GetSession(ctx context.Context, id uuid.UUID) (*auth.Sessions, error)
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*auth.Sessions, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
}
