package usersservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"usersservice/internal/domain/models"
	serviceerror "usersservice/internal/service"
	storageerror "usersservice/internal/storage"
	"usersservice/pkg/lib/logger/sl"

	"github.com/google/uuid"
)

type IUsersStorage interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type UsersService struct {
	log     *slog.Logger
	storage IUsersStorage
}

func New(log *slog.Logger, storage IUsersStorage) *UsersService {
	return &UsersService{
		log:     log,
		storage: storage,
	}
}

// GetUsers implements grpcapp.IUsersService.
func (u *UsersService) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "service.users.GetUsers"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	users, err := u.storage.GetUsers(ctx)
	if err != nil {
		log.Error("Error fetching users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// GetUserById implements grpcapp.IUsersService.
func (u *UsersService) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "service.users.GetUserById"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.storage.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerror.ErrNotFound))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
		}

		log.Error("Error fetching user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Insert implements grpcapp.IUsersService.
func (u *UsersService) Insert(ctx context.Context, userForInsert models.User) (models.User, error) {
	const op = "service.users.Insert"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	insertedUser, err := u.storage.Insert(ctx, userForInsert)
	if err != nil {
		if errors.Is(err, storageerror.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(storageerror.ErrAlreadyExists))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrAlreadyExists)
		}

		log.Error("Cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return insertedUser, nil
}

// Update implements grpcapp.IUsersService.
func (u *UsersService) Update(ctx context.Context, uid uuid.UUID, userForUpdate models.User) (models.User, error) {
	const op = "service.users.Update"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	updatedUser, err := u.storage.Update(ctx, uid, userForUpdate)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(storageerror.ErrNotFound))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
		}

		log.Error("Error updating user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedUser, nil
}

// Delete implements grpcapp.IUsersService.
func (u *UsersService) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "service.users.Delete"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	deletedUser, err := u.storage.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(storageerror.ErrNotFound))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
		}

		log.Error("Error deleting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return deletedUser, nil
}
