package authservice

import (
	"api-gateway/internal/domain/models"
	"api-gateway/pkg/lib/logger/sl"
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

type IAuthServer interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

type AuthService struct {
	log        *slog.Logger
	authServer IAuthServer
}

func New(log *slog.Logger, authServer IAuthServer) *AuthService {
	return &AuthService{
		log:        log,
		authServer: authServer,
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

	accessToken, refreshToken, err := a.authServer.Login(ctx, login, password)
	if err != nil {
		log.Error("Cannot login", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	_ = accessToken
	_ = refreshToken

	return accessToken, "", nil
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

	registeredUser, err := a.authServer.Register(ctx, userForRegister)
	if err != nil {
		log.Error("Cannot register", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return registeredUser, nil
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

	isAdmin, err := a.authServer.IsAdmin(ctx, uid)
	if err != nil {
		log.Error("Cannot check is an user admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
