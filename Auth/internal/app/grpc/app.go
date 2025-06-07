package grpcapp

import (
	"auth/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type IAuthService interface {
	Login(ctx context.Context, login string, password []byte) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

func New(log *slog.Logger, authService IAuthService, port int) *App {
	gRPCServer := grpc.NewServer()

	// authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
