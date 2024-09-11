package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/application/session"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

type AuthService struct {
	tokenService   inbound.TokenService
	sessionService *session.SessionService
}

func NewAuthService(tokenService inbound.TokenService, sessionService *session.SessionService) *AuthService {
	return &AuthService{tokenService, sessionService}
}

func (as *AuthService) CreateToken(
	user *user.User,
	durationStr string,
) (string, *auth.Payload, error) {
	token, payload, err := as.tokenService.CreateToken(user, durationStr)
	if err != nil {
		return "", nil, fmt.Errorf("auth service - create token: %w", err)
	}

	return token, payload, nil
}

func (as *AuthService) VerifyToken(token string) (*auth.Payload, error) {
	payload, err := as.tokenService.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("auth service - verify token: %w", err)
	}

	return payload, nil
}

func (as *AuthService) CreateSession(
	ctx context.Context,
	user *user.User,
	refreshToken string,
	refreshPayload *auth.Payload,
	userAgent string,
	clientIP string,
) (*auth.Sessions, error) {
	session, err := auth.NewSessions(
		user.ID(),
		refreshToken,
		refreshPayload,
		userAgent,
		clientIP,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid session parameters: %w", err)
	}

	err = as.sessionService.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return session, nil
}

func (as *AuthService) GetSession(ctx context.Context, sessionID uuid.UUID) (*auth.Sessions, error) {
	session, err := as.sessionService.GetSessionByUserID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}
