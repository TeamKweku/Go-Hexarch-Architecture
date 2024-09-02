package inbound

import (
	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

type TokenService interface {
	CreateToken(user *user.User, cfg config.Config) (string, *auth.Payload, error)
	VerifyToken(token string) (*auth.Payload, error)
}
