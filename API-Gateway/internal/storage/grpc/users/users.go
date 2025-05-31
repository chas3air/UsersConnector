package grpcstorage

import (
	"api-gateway/internal/domain/models"
	umprofiles "api-gateway/internal/domain/profiles/um"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"
	"strconv"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCUsersStorage struct {
	log  *slog.Logger
	conn *grpc.ClientConn
}

func New(log *slog.Logger, host string, port_s string) *GRPCUsersStorage {
	port, _ := strconv.Atoi(port_s)
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

func (u *GRPCUsersStorage) Close() {
	if err := u.conn.Close(); err != nil {
		panic(err)
	}
}

// GetUsers implements users.IUsersStorage.
func (s *GRPCUsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.grpc.users.GetUsers"
	log := s.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	c := umv1.NewUsersManagerClient(s.conn)
	res, err := c.GetUsers(ctx, nil)
	if err != nil {
		log.Warn("failed to get users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var resUsers = make([]models.User, 0, len(res.GetUsers()))
	for _, pbUser := range res.GetUsers() {
		user, err := umprofiles.ProtoUsrToUsr(pbUser)
		if err != nil {
			log.Warn("failed to convert proto user to model user", sl.Err(err))
			continue
		}
		resUsers = append(resUsers, user)
	}

	return resUsers, nil
}

// GetUserById implements users.IUsersStorage.
func (u *GRPCUsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements users.IUsersStorage.
func (u *GRPCUsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements users.IUsersStorage.
func (u *GRPCUsersStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements users.IUsersStorage.
func (u *GRPCUsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
