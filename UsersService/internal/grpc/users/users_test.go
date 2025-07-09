package usersgrpc_test

import (
	"context"
	"testing"
	"usersservice/internal/domain/models"
	"usersservice/internal/domain/profiles"
	usersgrpc "usersservice/internal/grpc/users"
	serviceerror "usersservice/internal/service"
	"usersservice/pkg/lib/logger"

	umv1 "github.com/chas3air/protos/gen/go/usersManager"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- Mock IUsersService ---

type MockUsersService struct {
	mock.Mock
}

func (m *MockUsersService) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUsersService) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersService) Insert(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersService) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	args := m.Called(ctx, uid, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersService) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

// --- Helpers ---

func newTestServer(t *testing.T, service *MockUsersService) *usersgrpc.ServerAPI {
	logger := logger.SetupLogger("local")
	return &usersgrpc.ServerAPI{
		Service: service,
		Log:     logger,
	}
}

// --- Tests ---

func TestGetUsers_Success(t *testing.T) {
	mockSvc := new(MockUsersService)
	users := []models.User{
		{Id: uuid.New(), Login: "user1"},
		{Id: uuid.New(), Login: "user2"},
	}
	mockSvc.On("GetUsers", mock.Anything).Return(users, nil)

	srv := newTestServer(t, mockSvc)
	resp, err := srv.GetUsers(context.Background(), &umv1.GetUsersRequest{})

	assert.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	transferedUser, _ := profiles.ProtoUsrToUsr(resp.GetUsers()[0])
	assert.Equal(t, users[0].Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestGetUsers_ContextCanceled(t *testing.T) {
	mockSvc := new(MockUsersService)
	srv := newTestServer(t, mockSvc)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := srv.GetUsers(ctx, &umv1.GetUsersRequest{})
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
}

func TestGetUserById_Success(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()
	user := models.User{Id: id, Login: "user1"}

	mockSvc.On("GetUserById", mock.Anything, id).Return(user, nil)

	srv := newTestServer(t, mockSvc)
	req := &umv1.GetUserByIdRequest{Id: id.String()}
	resp, err := srv.GetUserById(context.Background(), req)

	assert.NoError(t, err)

	transferedUser, _ := profiles.ProtoUsrToUsr(resp.User)
	assert.Equal(t, user.Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestGetUserById_InvalidUUID(t *testing.T) {
	mockSvc := new(MockUsersService)
	srv := newTestServer(t, mockSvc)

	req := &umv1.GetUserByIdRequest{Id: "invalid-uuid"}
	_, err := srv.GetUserById(context.Background(), req)

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetUserById_NotFound(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()

	mockSvc.On("GetUserById", mock.Anything, id).Return(models.User{}, serviceerror.ErrNotFound)

	srv := newTestServer(t, mockSvc)
	req := &umv1.GetUserByIdRequest{Id: id.String()}
	_, err := srv.GetUserById(context.Background(), req)

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockSvc.AssertExpectations(t)
}

func TestInsert_Success(t *testing.T) {
	mockSvc := new(MockUsersService)
	user := models.User{Id: uuid.New(), Login: "user1"}

	mockSvc.On("Insert", mock.Anything, user).Return(user, nil)

	srv := newTestServer(t, mockSvc)
	pbUser := profiles.UsrToProtoUsr(user)
	req := &umv1.InsertRequest{User: pbUser}

	resp, err := srv.Insert(context.Background(), req)
	assert.NoError(t, err)

	transferedUser, _ := profiles.ProtoUsrToUsr(resp.User)
	assert.Equal(t, user.Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestInsert_InvalidUser(t *testing.T) {
	mockSvc := new(MockUsersService)
	srv := newTestServer(t, mockSvc)

	req := &umv1.InsertRequest{User: &umv1.User{Id: "bad-uuid"}}
	_, err := srv.Insert(context.Background(), req)

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestInsert_AlreadyExists(t *testing.T) {
	mockSvc := new(MockUsersService)
	user := models.User{Id: uuid.New(), Login: "user1"}

	mockSvc.On("Insert", mock.Anything, user).Return(models.User{}, serviceerror.ErrAlreadyExists)

	srv := newTestServer(t, mockSvc)
	req := &umv1.InsertRequest{User: profiles.UsrToProtoUsr(user)}

	_, err := srv.Insert(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	mockSvc.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()
	user := models.User{Id: id, Login: "user1"}

	mockSvc.On("Update", mock.Anything, id, user).Return(user, nil)

	srv := newTestServer(t, mockSvc)
	req := &umv1.UpdateRequest{
		Id:   id.String(),
		User: profiles.UsrToProtoUsr(user),
	}

	resp, err := srv.Update(context.Background(), req)
	assert.NoError(t, err)

	transferedUser, _ := profiles.ProtoUsrToUsr(resp.User)
	assert.Equal(t, user.Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestUpdate_InvalidUUID(t *testing.T) {
	mockSvc := new(MockUsersService)
	srv := newTestServer(t, mockSvc)

	req := &umv1.UpdateRequest{Id: "bad-uuid"}
	_, err := srv.Update(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdate_NotFound(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()
	user := models.User{Id: id, Login: "user1"}

	mockSvc.On("Update", mock.Anything, id, user).Return(models.User{}, serviceerror.ErrNotFound)

	srv := newTestServer(t, mockSvc)
	req := &umv1.UpdateRequest{
		Id:   id.String(),
		User: profiles.UsrToProtoUsr(user),
	}
	_, err := srv.Update(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockSvc.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()
	user := models.User{Id: id, Login: "user1"}

	mockSvc.On("Delete", mock.Anything, id).Return(user, nil)

	srv := newTestServer(t, mockSvc)
	req := &umv1.DeleteRequest{Id: id.String()}

	resp, err := srv.Delete(context.Background(), req)
	assert.NoError(t, err)

	transferedUser, _ := profiles.ProtoUsrToUsr(resp.User)
	assert.Equal(t, user.Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestDelete_InvalidUUID(t *testing.T) {
	mockSvc := new(MockUsersService)
	srv := newTestServer(t, mockSvc)

	req := &umv1.DeleteRequest{Id: "bad-uuid"}
	_, err := srv.Delete(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDelete_NotFound(t *testing.T) {
	mockSvc := new(MockUsersService)
	id := uuid.New()

	mockSvc.On("Delete", mock.Anything, id).Return(models.User{}, serviceerror.ErrNotFound)

	srv := newTestServer(t, mockSvc)
	req := &umv1.DeleteRequest{Id: id.String()}
	_, err := srv.Delete(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockSvc.AssertExpectations(t)
}
