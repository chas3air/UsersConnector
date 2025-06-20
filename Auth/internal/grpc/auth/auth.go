package authgrpc

import (
	"auth/internal/domain/models"
	"context"
	"log/slog"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type IAuthService interface {
	Login(ctx context.Context, login string, password []byte) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	service IAuthService
	log     *slog.Logger
}

func Register(grpc *grpc.Server, service IAuthService, log *slog.Logger) {
	authv1.RegisterAuthServer(
		grpc,
		&serverAPI{
			service: service,
			log:     log,
		},
	)
}
