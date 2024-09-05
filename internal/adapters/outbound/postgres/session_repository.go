package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/postgres/sqlc"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
)

var _ outbound.SessionRepository = (*Client)(nil)

func (c *Client) CreateSession(ctx context.Context, session *auth.Sessions) error {
	arg := sqlc.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIP,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
	}

	_, err := c.queries.CreateSession(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (c *Client) GetSession(ctx context.Context, id uuid.UUID) (*auth.Sessions, error) {
	session, err := c.queries.GetSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &auth.Sessions{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}, nil
}

func (c *Client) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*auth.Sessions, error) {
	session, err := c.queries.GetSessionByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by user ID: %w", err)
	}

	return &auth.Sessions{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIp,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}, nil
}

func (c *Client) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := c.queries.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
