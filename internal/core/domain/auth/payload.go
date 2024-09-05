package auth

import (
	"time"

	"github.com/google/uuid"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

type Payload struct {
	ID        uuid.UUID `json:"id,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Role      string    `json:"role,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

// New Payload creates a new token payload with specific username and duration
func NewPayload(user *user.User, cfg config.Config) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, NewInvalidUUIDError("invalid UUID creation error")
	}

	durationStr := cfg.AccessTokenDuration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, NewInvalidDurationError(duration, "invalid token duration")
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    user.ID(),
		Role:      user.Role().String(),
		IssuedAt:  time.Now().UTC(),
		ExpiredAt: time.Now().UTC().Add(duration),
	}

	return payload, nil
}

// Valid checks if token payload is valid or not
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return NewExpiredTokenError(p.ExpiredAt)
	}

	return nil
}
