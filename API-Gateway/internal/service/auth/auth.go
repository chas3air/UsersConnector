package authservice

import (
	"api-gateway/internal/domain/models"
	"context"
	"fmt"
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

type IAuthStorage interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
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

// Login implements auth.IAuthService.
func (a *AuthService) Login(ctx context.Context, login string, password string) (string, string, error) {
	const op = "service.auth.Login"
	log := a.log.With(
		"op", op,
	)
	_ = log

	select {
	case <-ctx.Done():
		return "", "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	return "", "", nil
}

// Register implements auth.IAuthService.
func (a *AuthService) Register(ctx context.Context, userForRegister models.User) (models.User, error) {
	const op = "service.auth.Register"
	log := a.log.With(
		"op", op,
	)
	_ = log

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	return models.User{}, nil
}

// IsAdmin implements auth.IAuthService.
func (a *AuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	const op = "service.auth.IsAdmin"
	log := a.log.With(
		"op", op,
	)
	_ = log

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	return false, nil
}
