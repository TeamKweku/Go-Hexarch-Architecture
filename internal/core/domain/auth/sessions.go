package auth

import (
	"time"

	"github.com/google/uuid"
)

type Sessions struct {
	ID           uuid.UUID `json:"id,omitempty"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	ClientIP     string    `json:"client_ip,omitempty"`
	IsBlocked    bool      `json:"is_blocked,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// create a new session
func NewSessions(
	userID uuid.UUID,
	refreshToken string,
	refreshPayload *Payload,
	userAgent string,
	clientIP string,
) (*Sessions, error) {
	return &Sessions{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		CreatedAt:    time.Now().UTC(),
	}, nil
}
