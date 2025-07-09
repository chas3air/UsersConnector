package grpcusers_test

import (
	"context"
	"errors"
	"testing"

	"auth/internal/domain/models"
	umprofiles "auth/internal/profiles/um"
	grpcusers "auth/internal/storage/grpc/users"
	"auth/pkg/lib/logger"

	"log/slog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
)

// --- Mock UsersManagerClient ---

type MockUsersManagerClient struct {
	mock.Mock
}

func (m *MockUsersManagerClient) GetUsers(ctx context.Context, in *umv1.GetUsersRequest, opts ...interface{}) (*umv1.GetUsersResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*umv1.GetUsersResponse), args.Error(1)
}

func (m *MockUsersManagerClient) GetUserById(ctx context.Context, in *umv1.GetUserByIdRequest, opts ...interface{}) (*umv1.GetUserByIdResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*umv1.GetUserByIdResponse), args.Error(1)
}

func (m *MockUsersManagerClient) Insert(ctx context.Context, in *umv1.InsertRequest, opts ...interface{}) (*umv1.InsertResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*umv1.InsertResponse), args.Error(1)
}

// --- Mock GRPCUsersStorage with injected client ---

type MockGRPCUsersStorage struct {
	*grpcusers.GRPCUsersStorage
	mockClient *MockUsersManagerClient
}

func NewMockGRPCUsersStorage(log *slog.Logger, client *MockUsersManagerClient) *MockGRPCUsersStorage {
	return &MockGRPCUsersStorage{
		GRPCUsersStorage: &grpcusers.GRPCUsersStorage{Log: log},
		mockClient:       client,
	}
}

// Override methods to use mockClient instead of real grpc client

func (s *MockGRPCUsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.grpc.users.GetUsers"
	log := s.Log.With("op", op)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	res, err := s.mockClient.GetUsers(ctx, &umv1.GetUsersRequest{})
	if err != nil {
		log.Error("Failed to get users", slog.String("error", err.Error()))
		return nil, err
	}

	var resUsers []models.User
	for _, pbUser := range res.GetUsers() {
		u, err := umprofiles.ProtoUsrToUsr(pbUser)
		if err != nil {
			log.Warn("Error converting user", slog.String("error", err.Error()))
			continue
		}
		resUsers = append(resUsers, u)
	}
	return resUsers, nil
}

func (s *MockGRPCUsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.grpc.users.GetUserById"
	log := s.Log.With("op", op)

	select {
	case <-ctx.Done():
		return models.User{}, ctx.Err()
	default:
	}

	res, err := s.mockClient.GetUserById(ctx, &umv1.GetUserByIdRequest{Id: uid.String()})
	if err != nil {
		return models.User{}, err
	}

	user, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", slog.String("error", err.Error()))
		return models.User{}, err
	}
	return user, nil
}

func (s *MockGRPCUsersStorage) Insert(ctx context.Context, userForInsert models.User) (models.User, error) {
	const op = "storage.grpc.users.Insert"
	log := s.Log.With("op", op)

	select {
	case <-ctx.Done():
		return models.User{}, ctx.Err()
	default:
	}

	pbUser := umprofiles.UsrToProtoUsr(userForInsert)
	res, err := s.mockClient.Insert(ctx, &umv1.InsertRequest{User: pbUser})
	if err != nil {
		return models.User{}, err
	}

	insertedUser, err := umprofiles.ProtoUsrToUsr(res.GetUser())
	if err != nil {
		log.Error("Error converting user", slog.String("error", err.Error()))
		return models.User{}, err
	}
	return insertedUser, nil
}

// --- Tests ---

func TestGetUsers_Success(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	users := []*umv1.User{
		{Id: uuid.New().String(), Login: "user1"},
		{Id: uuid.New().String(), Login: "user2"},
	}
	mockClient.On("GetUsers", mock.Anything, &umv1.GetUsersRequest{}).
		Return(&umv1.GetUsersResponse{Users: users}, nil)

	got, err := storage.GetUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	mockClient.AssertExpectations(t)
}

func TestGetUsers_Error(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	mockClient.On("GetUsers", mock.Anything, &umv1.GetUsersRequest{}).
		Return((*umv1.GetUsersResponse)(nil), errors.New("some grpc error"))

	_, err := storage.GetUsers(context.Background())
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestGetUserById_Success(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	id := uuid.New()
	pbUser := &umv1.User{Id: id.String(), Login: "user1"}
	mockClient.On("GetUserById", mock.Anything, &umv1.GetUserByIdRequest{Id: id.String()}).
		Return(&umv1.GetUserByIdResponse{User: pbUser}, nil)

	got, err := storage.GetUserById(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "user1", got.Login)
	mockClient.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	id := uuid.New()
	mockClient.On("GetUserById", mock.Anything, &umv1.GetUserByIdRequest{Id: id.String()}).
		Return((*umv1.GetUserByIdResponse)(nil), status.Error(codes.NotFound, "user not found"))

	_, err := storage.GetUserById(context.Background(), id)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	mockClient.AssertExpectations(t)
}

func TestInsert_Success(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	user := models.User{Id: uuid.New(), Login: "user1"}
	pbUser := umprofiles.UsrToProtoUsr(user)
	mockClient.On("Insert", mock.Anything, &umv1.InsertRequest{User: pbUser}).
		Return(&umv1.InsertResponse{User: pbUser}, nil)

	got, err := storage.Insert(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, user.Login, got.Login)
	mockClient.AssertExpectations(t)
}

func TestInsert_AlreadyExists(t *testing.T) {
	log := logger.SetupLogger("local")
	mockClient := new(MockUsersManagerClient)
	storage := NewMockGRPCUsersStorage(log, mockClient)

	user := models.User{Id: uuid.New(), Login: "user1"}
	pbUser := umprofiles.UsrToProtoUsr(user)

	mockClient.On("Insert", mock.Anything, &umv1.InsertRequest{User: pbUser}).
		Return((*umv1.InsertResponse)(nil), status.Error(codes.AlreadyExists, "user already exists"))

	_, err := storage.Insert(context.Background(), user)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok, "error should be a gRPC status error")
	assert.Equal(t, codes.AlreadyExists, st.Code())

	mockClient.AssertExpectations(t)
}
