package usersservice

import (
	"api-gateway/internal/domain/models"
	storageerror "api-gateway/internal/storage"
	"api-gateway/pkg/lib/logger/sl"
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

// GetUsers implements IUsersStorage.
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
		log.Error("Cannot fetxh users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// GetUserById implements IUsersStorage.
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
			log.Warn("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Cannot fetch user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Insert implements IUsersStorage.
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
			log.Warn("User already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Cannot insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return insertedUser, nil
}

// Update implements IUsersStorage.
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
			log.Warn("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Cannot update user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedUser, nil
}

// Delete implements IUsersStorage.
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
			log.Warn("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("Cannot delete user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return deletedUser, nil
}
