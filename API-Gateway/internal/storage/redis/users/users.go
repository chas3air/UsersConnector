package userscashstorage

import (
	"api-gateway/internal/domain/models"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UsersCashStorage struct {
	log            *slog.Logger
	rds            *redis.Client
	expirationTime int
}

func New(log *slog.Logger, host string, port int, expirationTime int) *UsersCashStorage {
	rds := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	if err := rds.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	return &UsersCashStorage{
		log:            log,
		rds:            rds,
		expirationTime: expirationTime,
	}
}

func (u *UsersCashStorage) Close() {
	u.rds.Close()
}

// Get implements userscashservice.UsersCashStorage.
func (u *UsersCashStorage) Get(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.redis.users.Get"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	userFromRedis, err := u.rds.HGetAll(
		ctx,
		fmt.Sprintf("user:%s", id.String()),
	).Result()
	if err != nil {
		log.Info("User not found in redis", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userForReturn := mapToUser(userFromRedis)
	return userForReturn, nil
}

// Set implements userscashservice.UsersCashStorage.
func (u *UsersCashStorage) Set(ctx context.Context, user models.User) error {
	const op = "storage.redis.users.Set"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := u.rds.HSet(
		ctx,
		fmt.Sprintf("user:%s", user.Id.String()),
		map[string]string{
			"id":       user.Id.String(),
			"login":    user.Login,
			"password": string(user.Password),
			"role":     user.Role,
		},
	).Result()
	if err != nil {
		log.Warn("Cannot insert user to redis", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.rds.Expire(
		ctx,
		fmt.Sprintf("user:%s", user.Id.String()),
		time.Duration(u.expirationTime),
	).Result()
	if err != nil {
		log.Warn("Cannot set expiration time for user"+user.Id.String(), sl.Err(err))
	}

	return nil
}

// Del implements userscashservice.UsersCashStorage.
func (u *UsersCashStorage) Del(ctx context.Context, id uuid.UUID) error {
	const op = "storage.redis.users.Del"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := u.rds.HDel(
		ctx,
		fmt.Sprintf("user:%s", id.String()),
	).Result()
	if err != nil {
		log.Warn("Cannot delete user:"+id.String()+" from cash", sl.Err(err))
	}

	return nil

}

func mapToUser(mappedUser map[string]string) models.User {
	id, _ := uuid.Parse(mappedUser["id"])
	return models.User{
		Id:       id,
		Login:    mappedUser["login"],
		Password: mappedUser["password"],
		Role:     mappedUser["role"],
	}
}
