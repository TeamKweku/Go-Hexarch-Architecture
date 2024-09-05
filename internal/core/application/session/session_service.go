package session

import (
	"context"

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
	return s.sessionRepo.CreateSession(ctx, session)
}

func (s *SessionService) GetSession(ctx context.Context, id uuid.UUID) (*auth.Sessions, error) {
	return s.sessionRepo.GetSession(ctx, id)
}

func (s *SessionService) GetSessionByUserID(ctx context.Context, userID uuid.UUID) (*auth.Sessions, error) {
	return s.sessionRepo.GetSessionByUserID(ctx, userID)
}

func (s *SessionService) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return s.sessionRepo.DeleteSession(ctx, id)
}
