package users

import (
	"api-gateway/internal/domain/models"
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type IUsersService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type UsersHandler struct {
	log     *slog.Logger
	service IUsersService
}

func New(log *slog.Logger, service IUsersService) *UsersHandler {
	return &UsersHandler{
		log:     log,
		service: service,
	}
}

func (u *UsersHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request)    {}
func (u *UsersHandler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {}
func (u *UsersHandler) InsertHandler(w http.ResponseWriter, r *http.Request)      {}
func (u *UsersHandler) UpdateHandler(w http.ResponseWriter, r *http.Request)      {}
func (u *UsersHandler) DeleteHandler(w http.ResponseWriter, r *http.Request)      {}
