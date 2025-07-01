package authgrpc

import (
	"auth/internal/domain/models"
	amprofiles "auth/internal/profiles/am"
	serviceerrors "auth/internal/service"
	"auth/pkg/lib/logger/sl"
	"context"
	"errors"
	"log/slog"

	authv1 "github.com/chas3air/protos/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthService interface {
	Login(ctx context.Context, email string, password string) (string, string, error)
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

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	const op = "grpc.auth.Login"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return nil, status.Error(codes.DeadlineExceeded, "context is over")
	default:
	}

	accessToken, refreshToken, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		log.Error("Cannot generate token", sl.Err(err))
		return nil, status.Error(codes.Internal, "Cannot generate token")
	}

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	const op = "grpc.auth.Register"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return nil, status.Error(codes.DeadlineExceeded, "context is over")
	default:
	}

	userforRegister, err := amprofiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		log.Error("Invalid argument", sl.Err(serviceerrors.ErrInvalidArgument))
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	createdUser, err := s.service.Register(ctx, userforRegister)
	if err != nil {
		switch {
		case errors.Is(err, serviceerrors.ErrDeadlineExceeded):
			log.Warn("Deadline exceeded", sl.Err(serviceerrors.ErrDeadlineExceeded))
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")

		case errors.Is(err, serviceerrors.ErrInvalidArgument):
			log.Warn("Invalid argument", sl.Err(serviceerrors.ErrInvalidArgument))
			return nil, status.Error(codes.InvalidArgument, "invalid argument")

		case errors.Is(err, serviceerrors.ErrAlreadyExists):
			log.Warn("User already exists", sl.Err(serviceerrors.ErrAlreadyExists))
			return nil, status.Error(codes.AlreadyExists, "user already exists")

		case errors.Is(err, serviceerrors.ErrNotFound):
			log.Warn("User not found", sl.Err(serviceerrors.ErrNotFound))
			return nil, status.Error(codes.NotFound, "user not found")

		default:
			log.Error("Cannot register user", sl.Err(err))
			return nil, status.Error(codes.Internal, "cannot register user")
		}
	}

	return &authv1.RegisterResponse{
		User: amprofiles.UsrToProtoUsr(createdUser),
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	const op = "grpc.auth.IsAdmin"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Info("context is over")
		return nil, status.Error(codes.DeadlineExceeded, "context is over")
	default:
	}

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Error("Invalid argument", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	isAdmin, err := s.service.IsAdmin(ctx, id)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrDeadlineExceeded) {
			log.Warn("Deadline exceeded", sl.Err(serviceerrors.ErrDeadlineExceeded))
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		} else if errors.Is(err, serviceerrors.ErrInvalidArgument) {
			log.Warn("Invalid argument", sl.Err(serviceerrors.ErrInvalidArgument))
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		} else if errors.Is(err, serviceerrors.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerrors.ErrNotFound))
			return nil, status.Error(codes.NotFound, "User not found")
		} else {
			log.Error("Cannot retrieve user", sl.Err(err))
			return nil, status.Error(codes.Internal, "Cannot retrieve user")
		}
	}

	return &authv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
