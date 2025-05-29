package app

import (
	"context"
	"log/slog"
	grpcapp "usersservice/internal/app/grpc"
	"usersservice/internal/domain/models"
	usersservice "usersservice/internal/service/users"

	"github.com/google/uuid"
)

type App struct {
	GRPCServer *grpcapp.App
}

type IUsersStorage interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

func New(log *slog.Logger, port int, storage IUsersStorage) *App {
	usersService := usersservice.New(log, storage)
	grpcapp := grpcapp.New(log, usersService, port)

	return &App{
		GRPCServer: grpcapp,
	}
}
