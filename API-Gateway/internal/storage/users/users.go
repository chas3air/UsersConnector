package users

import (
	"api-gateway/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UsersStorage struct {
	log  *slog.Logger
	host string
	port int
}

func New(log *slog.Logger, host string, port int) *UsersStorage {
	return &UsersStorage{
		log:  log,
		host: host,
		port: port,
	}
}

// GetUsers implements users.IUsersStorage.
func (u *UsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements users.IUsersStorage.
func (u *UsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements users.IUsersStorage.
func (u *UsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements users.IUsersStorage.
func (u *UsersStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements users.IUsersStorage.
func (u *UsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
