package app

import (
	grpcapp "auth/internal/app/grpc"
	"auth/internal/domain/models"
	authservice "auth/internal/service/auth"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type App struct {
	GRPCServer *grpcapp.App
}

type IUsersStorage interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
}

func New(log *slog.Logger, port int, storage IUsersStorage) *App {
	authService := authservice.New(log, storage)
	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
