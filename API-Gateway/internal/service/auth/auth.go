package authservice

import (
	"api-gateway/internal/domain/models"
	serviceerror "api-gateway/internal/service"
	storageerror "api-gateway/internal/storage"
	"api-gateway/pkg/lib/logger/sl"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

const isAdminTitleRole = "admin"

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

// Login implements auth.IAuthService.
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
		log.Error("Cannot fetch users", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	var bufUser models.User

	for _, user := range users {
		if login == user.Login && bytes.Equal(password, user.Password) {
			bufUser = user
		}
	}
	if bufUser.Login == "" {
		log.Error("User not found", sl.Err(serviceerror.ErrNotFound))
		return "", "", fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
	}

	// написать функцию генерации jwt-токена
	// использовать для этого bufUser

	return bufUser.Id.String(), bufUser.Login, nil
}

// Register implements auth.IAuthService.
func (a *AuthService) Register(ctx context.Context, userForRegister models.User) (models.User, error) {
	const op = "service.auth.Register"
	log := a.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	registeredUser, err := a.storage.Insert(ctx, userForRegister)
	if err != nil {
		if errors.Is(err, storageerror.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrAlreadyExists)
		}

		log.Error("Cannot registed user", sl.Err(err))
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

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := a.storage.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
		}

		log.Error("Cannot fetch user by id", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if user.Role != isAdminTitleRole {
		log.Info("User's role is not admin")
		return false, nil
	} else {
		log.Info("User's role is not admin")
		return true, nil
	}
}
