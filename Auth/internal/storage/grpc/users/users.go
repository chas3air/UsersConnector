package grpcusers

import (
	"auth/internal/domain/models"
	umprofiles "auth/internal/profiles/um"
	storageerrors "auth/internal/storage"
	"auth/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GRPCUsersStorage struct {
	Log  *slog.Logger
	conn *grpc.ClientConn
}

func New(log *slog.Logger, host string, port int) *GRPCUsersStorage {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		panic(err)
	}

	return &GRPCUsersStorage{
		Log:  log,
		conn: conn,
	}
}

func (s *GRPCUsersStorage) Close() {
	if err := s.conn.Close(); err != nil {
		panic(err)
	}
}

// GetUsers implements users.IUsersStorage.
func (s *GRPCUsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.grpc.users.GetUsers"
	log := s.Log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.GetUsers(ctx, nil)
	if err != nil {
		log.Error("Failed to get users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var resUsers = make([]models.User, 0, len(res.GetUsers()))
	for _, pbUser := range res.GetUsers() {
		convertedUser, err := umprofiles.ProtoUsrToUsr(pbUser)
		if err != nil {
			log.Warn("Error converting user", sl.Err(err))
			continue
		}

		resUsers = append(resUsers, convertedUser)
	}

	return resUsers, nil
}

// GetUserById implements users.IUsersStorage.
func (s *GRPCUsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.grpc.users.GetUserById"
	log := s.Log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.GetUserById(ctx, &umv1.GetUserByIdRequest{
		Id: uid.String(),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				log.Warn("Deadline exceeded", sl.Err(storageerrors.ErrDeadlineExceeded))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrDeadlineExceeded)
			case codes.InvalidArgument:
				log.Warn("Invalid argument", sl.Err(storageerrors.ErrInvalidArgument))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrInvalidArgument)
			case codes.NotFound:
				log.Warn("User not found", sl.Err(storageerrors.ErrNotFound))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
			default:
				log.Error("Cannot retrieve user by id", sl.Err(err))
				return models.User{}, fmt.Errorf("%s; %w", op, err)
			}
		}
	}

	convertedUser, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return convertedUser, nil
}

// Insert implements users.IUsersStorage.
func (s *GRPCUsersStorage) Insert(ctx context.Context, userForInsert models.User) (models.User, error) {
	const op = "storage.grpc.users.Insert"
	log := s.Log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.Insert(ctx, &umv1.InsertRequest{
		User: umprofiles.UsrToProtoUsr(userForInsert),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				log.Warn("Deadline exceeded", sl.Err(storageerrors.ErrDeadlineExceeded))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrDeadlineExceeded)
			case codes.InvalidArgument:
				log.Warn("Invalid argument", sl.Err(storageerrors.ErrInvalidArgument))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrInvalidArgument)
			case codes.AlreadyExists:
				log.Warn("User already exists", sl.Err(storageerrors.ErrAlreadyExists))
				return models.User{}, fmt.Errorf("%s: %w", op, storageerrors.ErrAlreadyExists)
			default:
				log.Error("Cannot retrieve user by id", sl.Err(err))
				return models.User{}, fmt.Errorf("%s; %w", op, err)
			}
		}
	}

	insertedUser, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", sl.Err(err))
	}

	return insertedUser, nil
}
