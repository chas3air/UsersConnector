package authhandler

import (
	"api-gateway/internal/domain/models"
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type IAuthService interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

type AuthHandler struct {
	log     *slog.Logger
	service IAuthService
}

func New(log *slog.Logger, service IAuthService) *AuthHandler {
	return &AuthHandler{
		log:     log,
		service: service,
	}
}

func (AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request)        {}
func (AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request)     {}
func (AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {}
func (AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request)       {}
