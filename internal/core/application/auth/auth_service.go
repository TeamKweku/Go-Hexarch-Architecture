package auth

import (
	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

type AuthService struct {
	tokenService inbound.TokenService
}

func NewAuthService(tokenService inbound.TokenService) *AuthService {
	return &AuthService{tokenService}
}

func (as *AuthService) CreateToken(
	user *user.User,
	cfg config.Config,
) (string, *auth.Payload, error) {
	return as.tokenService.CreateToken(user, cfg)
}

func (as *AuthService) VerifyToken(token string) (*auth.Payload, error) {
	return as.tokenService.VerifyToken(token)
}
