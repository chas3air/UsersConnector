package grpcusers

import (
	"auth/internal/domain/models"
	"auth/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

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
	panic("unimplemented")
}

// GetUserById implements users.IUsersStorage.
func (s *GRPCUsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements users.IUsersStorage.
func (s *GRPCUsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements users.IUsersStorage.
func (s *GRPCUsersStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements users.IUsersStorage.
func (s *GRPCUsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
