package grpcapp

import (
	"auth/internal/domain/models"
	authgrpc "auth/internal/grpc/auth"
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type IAuthService interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

func New(log *slog.Logger, authService IAuthService, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc.Run"
	log := a.log.With(
		"op", op,
	)

	l, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", a.port),
	)
	if err != nil {
		return err
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	log.Info("app started")

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
