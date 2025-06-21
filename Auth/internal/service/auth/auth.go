package authservice

import (
	"auth/internal/domain/models"
	"auth/pkg/lib/logger/sl"
	"bytes"
	"context"
	"errors"
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
	const op = "service.auth.Login"
	log := a.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return "", "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := a.storage.GetUsers(ctx)
	if err != nil {
		log.Error("Failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	var loggedUser = models.User{Login: "undef"}
	for _, user := range users {
		if user.Login == login && bytes.Equal(user.Password, password) {
			loggedUser = user
		}
	}
	if loggedUser.Login == "undef" {
		log.Error("User doesn't exists", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, errors.New("user doesn't exists"))
	}

	// генерация токенов
	accessToken, refreshToken, err := "", "", nil
	return accessToken, refreshToken, nil
}

// Register implements grpcapp.IAuthService.
func (a *AuthService) Register(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// IsAdmin implements grpcapp.IAuthService.
func (a *AuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	panic("unimplemented")
}
