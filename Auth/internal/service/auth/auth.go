package authservice

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	serviceerrors "auth/internal/service"
	storageerrors "auth/internal/storage"
	"auth/pkg/lib/logger/sl"
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
func (a *AuthService) Login(ctx context.Context, login string, password string) (string, string, error) {
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
		if user.Login == login && user.Password == password {
			loggedUser = user
		}
	}
	if loggedUser.Login == "undef" {
		log.Error("User doesn't exists")
		return "", "", fmt.Errorf("%s: %w", op, errors.New("user doesn't exists"))
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(loggedUser)
	if err != nil {
		log.Error("Failed to generate tokens", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

// Register implements grpcapp.IAuthService.
func (a *AuthService) Register(ctx context.Context, userForCheck models.User) (models.User, error) {
	const op = "service.auth.Register"
	log := a.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := a.storage.GetUsers(ctx)
	if err != nil {
		log.Error("Cannot fetching users", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var appUser models.User
	for _, user := range users {
		if userForCheck.Login == user.Login && userForCheck.Password == user.Password {
			appUser = user
		}
	}

	if appUser.Login != "" {
		log.Warn("User already exists", sl.Err(serviceerrors.ErrAlreadyExists))
		return models.User{}, fmt.Errorf("%s: %w", op, serviceerrors.ErrAlreadyExists)
	}

	insertedUser, err := a.storage.Insert(ctx, userForCheck)
	if err != nil {
		log.Error("Cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return insertedUser, nil
}

// IsAdmin implements grpcapp.IAuthService.
func (a *AuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	const op = "service.auth.IsAdmind"
	log := a.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := a.storage.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storageerrors.ErrDeadlineExceeded) {
			log.Warn("Deadline exceeded", sl.Err(serviceerrors.ErrDeadlineExceeded))
			return false, fmt.Errorf("%s: %w", op, serviceerrors.ErrDeadlineExceeded)
		} else if errors.Is(err, storageerrors.ErrInvalidArgument) {
			log.Warn("Invalid argument", sl.Err(serviceerrors.ErrInvalidArgument))
			return false, fmt.Errorf("%s: %w", op, serviceerrors.ErrInvalidArgument)
		} else if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerrors.ErrNotFound))
			return false, fmt.Errorf("%s: %w", op, serviceerrors.ErrNotFound)
		} else {
			log.Error("Cannot retrieve user", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, err)
		}
	}

	return user.Role == "admin", nil
}
