package usersmongostorage

import (
	"context"
	"log/slog"
	"usersservice/internal/domain/models"

	"github.com/google/uuid"
)

type UsersMongoStorage struct {
	log *slog.Logger
}

func New(log *slog.Logger) *UsersMongoStorage {
	return &UsersMongoStorage{
		log: log,
	}
}

func (u *UsersMongoStorage) Close() {
}

// GetUsers implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements usersservice.IUsersStorage.
func (u *UsersMongoStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
