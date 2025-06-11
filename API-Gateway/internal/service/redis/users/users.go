package userscashservice

import (
	"api-gateway/internal/domain/models"
	serviceerror "api-gateway/internal/service"
	storageerror "api-gateway/internal/storage"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type UsersCashStorage interface {
	Get(context.Context, uuid.UUID) (models.User, error)
	Set(context.Context, models.User) error
	Del(context.Context, uuid.UUID) error
}

type UsersCashService struct {
	log     *slog.Logger
	storage UsersCashStorage
}

func New(log *slog.Logger, storage UsersCashStorage) *UsersCashService {
	return &UsersCashService{
		log:     log,
		storage: storage,
	}
}

func (u *UsersCashService) Get(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "service.redis.users.Get"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	userFromCash, err := u.storage.Get(ctx, id)
	if err != nil {
		if err == storageerror.ErrNotFound {
			log.Warn("not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, serviceerror.ErrNotFound)
		}

		log.Error("Error fetching data from cash", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userFromCash, nil
}

func (u *UsersCashService) Set(ctx context.Context, user models.User) error {
	const op = "service.redis.users.Set"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	if err := u.storage.Set(ctx, user); err != nil {
		log.Error("Error set user to cash", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UsersCashService) Del(ctx context.Context, id uuid.UUID) error {
	const op = "service.redis.users.Del"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	if err := u.storage.Del(ctx, id); err != nil {
		log.Error("Erro deleting user from cash", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
