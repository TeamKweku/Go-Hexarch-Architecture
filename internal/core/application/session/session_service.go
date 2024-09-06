package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
)

type SessionService struct {
	sessionRepo outbound.SessionRepository
}

func NewSessionService(sessionRepo outbound.SessionRepository) *SessionService {
	return &SessionService{sessionRepo}
}

func (s *SessionService) CreateSession(ctx context.Context, session *auth.Sessions) error {
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (s *SessionService) GetSession(ctx context.Context, id uuid.UUID) (*auth.Sessions, error) {
	session, err := s.sessionRepo.GetSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (s *SessionService) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*auth.Sessions, error) {
	session, err := s.sessionRepo.GetSessionByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, id uuid.UUID) error {
	if err := s.sessionRepo.DeleteSession(ctx, id); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
