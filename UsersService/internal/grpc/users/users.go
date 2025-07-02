package usersgrpc

import (
	"context"
	"errors"
	"log/slog"
	"usersservice/internal/domain/models"
	"usersservice/internal/domain/profiles"
	serviceerror "usersservice/internal/service"
	"usersservice/pkg/lib/logger/sl"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUsersService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type serverAPI struct {
	umv1.UnimplementedUsersManagerServer
	service IUsersService
	log     *slog.Logger
}

func Register(grpc *grpc.Server, service IUsersService, log *slog.Logger) {
	umv1.RegisterUsersManagerServer(
		grpc,
		&serverAPI{
			service: service,
			log:     log,
		},
	)
}

// GetUsers implements umv1.UsersManagerServer.
func (s *serverAPI) GetUsers(ctx context.Context, req *umv1.GetUsersRequest) (*umv1.GetUsersResponse, error) {
	const op = "grpc.users.GetUsers"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("Request time out")
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	users, err := s.service.GetUsers(ctx)
	if err != nil {
		log.Error("Error fetching users", sl.Err(err))
		return nil, status.Error(codes.Internal, "error fetching users")
	}

	var responseUsers = make([]*umv1.User, 0, len(users))
	for _, user := range users {
		profiledUser := profiles.UsrToProtoUsr(user)
		responseUsers = append(responseUsers, profiledUser)
	}

	return &umv1.GetUsersResponse{
		Users: responseUsers,
	}, nil
}

// GetUserById implements umv1.UsersManagerServer.
func (s *serverAPI) GetUserById(ctx context.Context, req *umv1.GetUserByIdRequest) (*umv1.GetUserByIdResponse, error) {
	const op = "grpc.users.GetUserById"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("Request time out")
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	uid, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("Cannot parse request uid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid uid")
	}

	user, err := s.service.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerror.ErrNotFound))
			return nil, status.Error(codes.NotFound, "user not found")
		}

		log.Error("Error fetching user by id", sl.Err(err))
		return nil, status.Error(codes.Internal, "error fetching user by id")
	}

	return &umv1.GetUserByIdResponse{
		User: profiles.UsrToProtoUsr(user),
	}, nil
}

// Insert implements umv1.UsersManagerServer.
func (s *serverAPI) Insert(ctx context.Context, req *umv1.InsertRequest) (*umv1.InsertResponse, error) {
	const op = "grpc.users.Insert"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("Request time out")
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	userForInsert, err := profiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		log.Error("Error parse pb_user to user", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid user")
	}

	insertedUser, err := s.service.Insert(ctx, userForInsert)
	if err != nil {
		if errors.Is(err, serviceerror.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(serviceerror.ErrAlreadyExists))
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Error("Error inserting user", sl.Err(err))
		return nil, status.Error(codes.Internal, "error inserting user")
	}

	return &umv1.InsertResponse{
		User: profiles.UsrToProtoUsr(insertedUser),
	}, nil
}

// Update implements umv1.UsersManagerServer.
func (s *serverAPI) Update(ctx context.Context, req *umv1.UpdateRequest) (*umv1.UpdateResponse, error) {
	const op = "grpc.users.Update"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("Request time out")
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	uid, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("Cannot parse request uid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid uid")
	}

	userForUpdate, err := profiles.ProtoUsrToUsr(req.GetUser())
	if err != nil {
		log.Error("Error parse pb_user to user", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid user")
	}

	updatedUser, err := s.service.Update(ctx, uid, userForUpdate)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerror.ErrNotFound))
			return nil, status.Error(codes.NotFound, "user not found")
		}

		log.Error("Error updating user", sl.Err(err))
		return nil, status.Error(codes.Internal, "error updating user")
	}

	return &umv1.UpdateResponse{
		User: profiles.UsrToProtoUsr(updatedUser),
	}, nil
}

// Delete implements umv1.UsersManagerServer.
func (s *serverAPI) Delete(ctx context.Context, req *umv1.DeleteRequest) (*umv1.DeleteResponse, error) {
	const op = "grpc.users.Delete"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		log.Error("Request time out")
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	uid, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Error("Cannot parse request uid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid uid")
	}

	deletedUser, err := s.service.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(serviceerror.ErrNotFound))
			return nil, status.Error(codes.NotFound, "user not found")
		}

		log.Error("Error deleting user", sl.Err(err))
		return nil, status.Error(codes.Internal, "error deleting user")
	}

	return &umv1.DeleteResponse{
		User: profiles.UsrToProtoUsr(deletedUser),
	}, nil
}
