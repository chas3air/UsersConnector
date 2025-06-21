package grpcusers

import (
	"auth/internal/domain/models"
	umprofiles "auth/internal/profiles/um"
	"auth/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCUsersStorage struct {
	log  *slog.Logger
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
		log:  log,
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
	log := s.log.With(
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
	log := s.log.With(
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
		log.Error("Failed to get user by id", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
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
	log := s.log.With(
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
		log.Error("Failed to insert user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	insertedUser, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", sl.Err(err))
	}

	return insertedUser, nil
}

// Update implements users.IUsersStorage.
func (s *GRPCUsersStorage) Update(ctx context.Context, uid uuid.UUID, userForUpdate models.User) (models.User, error) {
	const op = "storage.grpc.users.Update"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.Update(ctx, &umv1.UpdateRequest{
		Id:   uid.String(),
		User: umprofiles.UsrToProtoUsr(userForUpdate),
	})
	if err != nil {
		log.Error("Failed to update user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	updatedUser, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedUser, nil
}

// Delete implements users.IUsersStorage.
func (s *GRPCUsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.grpc.users.Delete"
	log := s.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.Delete(ctx, &umv1.DeleteRequest{
		Id: uid.String(),
	})
	if err != nil {
		log.Error("Failed to delete user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	deletedUSer, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return deletedUSer, nil
}
