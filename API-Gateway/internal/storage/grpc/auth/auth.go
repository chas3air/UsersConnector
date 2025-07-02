package grpcauthserver

import (
	"api-gateway/internal/domain/models"
	asprofiles "api-gateway/internal/domain/profiles/as"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCAuthServer struct {
	log  *slog.Logger
	conn *grpc.ClientConn
}

func New(log *slog.Logger, host string, port int) *GRPCAuthServer {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		panic(err)
	}

	return &GRPCAuthServer{
		log:  log,
		conn: conn,
	}
}

func (u *GRPCAuthServer) Close() {
	if err := u.conn.Close(); err != nil {
		panic(err)
	}
}

// Login implements authservice.IAuthStorage.
func (u *GRPCAuthServer) Login(ctx context.Context, login string, password string) (string, string, error) {
	const op = "storage.grpc.auth.Login"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over", sl.Err(ctx.Err()))
		return "", "", fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := authv1.NewAuthClient(u.conn)
	res, err := c.Login(
		ctx,
		&authv1.LoginRequest{
			Login:    login,
			Password: password,
		},
	)
	if err != nil {
		log.Error("Cannot login user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return res.AccessToken, res.RefreshToken, nil
}

// Register implements authservice.IAuthStorage.
func (u *GRPCAuthServer) Register(ctx context.Context, userForRegister models.User) (models.User, error) {
	const op = "storage.grpc.auth.Register"
	log := u.log.With(
		"op", op,
	)
	select {
	case <-ctx.Done():
		log.Info("context is over")
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := authv1.NewAuthClient(u.conn)
	res, err := c.Register(ctx,
		&authv1.RegisterRequest{
			User: asprofiles.UsrToProtoUsr(userForRegister),
		})
	if err != nil {
		log.Error("Cannot register user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return asprofiles.ProtoUsrToUsr(res.GetUser())
}

// IsAdmin implements authservice.IAuthStorage.
func (u *GRPCAuthServer) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	const op = "storage.grpc.auth.IsAdmin"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return false, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := authv1.NewAuthClient(u.conn)
	res, err := c.IsAdmin(ctx,
		&authv1.IsAdminRequest{
			UserId: uid.String(),
		},
	)
	if err != nil {
		log.Error("Cannot check is an user admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res.IsAdmin, nil
}
