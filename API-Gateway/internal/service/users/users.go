package users

import (
	"api-gateway/internal/domain/models"
	"context"
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
	service IUsersStorage
}

func New(log *slog.Logger, service IUsersStorage) *UsersService {
	return &UsersService{
		log:     log,
		service: service,
	}
}

// GetUsers implements IUsersStorage.
func (u *UsersService) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements IUsersStorage.
func (u *UsersService) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements IUsersStorage.
func (u *UsersService) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements IUsersStorage.
func (u *UsersService) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements IUsersStorage.
func (u *UsersService) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
