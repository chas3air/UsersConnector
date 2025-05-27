package auth

import (
	"api-gateway/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type IUsersStorage interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type AuthService struct {
	log     *slog.Logger
	service IUsersStorage
}

func New(log *slog.Logger, service IUsersStorage) *AuthService {
	return &AuthService{
		log:     log,
		service: service,
	}
}

// Login implements auth.IAuthService.
func (a *AuthService) Login(ctx context.Context, login string, password string) (string, string, error) {
	panic("unimplemented")
}

// Register implements auth.IAuthService.
func (a *AuthService) Register(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// IsAdmin implements auth.IAuthService.
func (a *AuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	panic("unimplemented")
}
