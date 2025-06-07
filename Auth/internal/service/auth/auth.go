package authservice

import (
	"auth/internal/domain/models"
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
	storage IUsersStorage
}

func New(log *slog.Logger, storage IUsersStorage) *AuthService {
	return &AuthService{
		log:     log,
		storage: storage,
	}
}

// Login implements grpcapp.IAuthService.
func (a *AuthService) Login(ctx context.Context, login string, password []byte) (string, string, error) {
	panic("unimplemented")
}

// Register implements grpcapp.IAuthService.
func (a *AuthService) Register(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// IsAdmin implements grpcapp.IAuthService.
func (a *AuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	panic("unimplemented")
}
