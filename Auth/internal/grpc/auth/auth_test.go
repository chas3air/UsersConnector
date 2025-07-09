package authgrpc_test

import (
	"context"
	"errors"
	"testing"

	"auth/internal/domain/models"
	authgrpc "auth/internal/grpc/auth"
	amprofiles "auth/internal/profiles/am"
	serviceerrors "auth/internal/service"
	"auth/pkg/lib/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/chas3air/protos/gen/go/auth"
)

// --- Mock IAuthService ---

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) Register(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthService) IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	args := m.Called(ctx, uid)
	return args.Bool(0), args.Error(1)
}

// --- Helpers ---

func newTestServer(t *testing.T, service *MockAuthService) *authgrpc.ServerAPI {
	logger := logger.SetupLogger("local")
	return &authgrpc.ServerAPI{
		Service: service,
		Log:     logger,
	}
}

// --- Tests ---

func TestLogin_ContextDone(t *testing.T) {
	mockSvc := new(MockAuthService)
	srv := newTestServer(t, mockSvc)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &authv1.LoginRequest{Email: "email@test.com", Password: "password"}
	_, err := srv.Login(ctx, req)

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
}

func TestLogin_Error(t *testing.T) {
	mockSvc := new(MockAuthService)
	mockSvc.On("Login", mock.Anything, "email@test.com", "password").Return("", "", errors.New("some error"))

	srv := newTestServer(t, mockSvc)
	req := &authv1.LoginRequest{Email: "email@test.com", Password: "password"}

	_, err := srv.Login(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockSvc.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	mockSvc := new(MockAuthService)
	user := models.User{Login: "user1", Password: "pass"}
	pbUser := amprofiles.UsrToProtoUsr(user)

	mockSvc.On("Register", mock.Anything, user).Return(user, nil)

	srv := newTestServer(t, mockSvc)
	req := &authv1.RegisterRequest{User: pbUser}

	resp, err := srv.Register(context.Background(), req)
	assert.NoError(t, err)
	transferedUser, _ := amprofiles.ProtoUsrToUsr(resp.User)
	assert.Equal(t, user.Login, transferedUser.Login)
	mockSvc.AssertExpectations(t)
}

func TestRegister_InvalidProtoUser(t *testing.T) {
	mockSvc := new(MockAuthService)
	srv := newTestServer(t, mockSvc)

	// Passing nil user to cause ProtoUsrToUsr error
	req := &authv1.RegisterRequest{User: nil}

	_, err := srv.Register(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestRegister_Errors(t *testing.T) {
	mockSvc := new(MockAuthService)
	user := models.User{Login: "user1", Password: "pass"}

	tests := []struct {
		name       string
		serviceErr error
		wantCode   codes.Code
	}{
		{"DeadlineExceeded", serviceerrors.ErrDeadlineExceeded, codes.DeadlineExceeded},
		{"InvalidArgument", serviceerrors.ErrInvalidArgument, codes.InvalidArgument},
		{"AlreadyExists", serviceerrors.ErrAlreadyExists, codes.AlreadyExists},
		{"NotFound", serviceerrors.ErrNotFound, codes.NotFound},
		{"OtherError", errors.New("other"), codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.ExpectedCalls = nil
			mockSvc.On("Register", mock.Anything, user).Return(models.User{}, tt.serviceErr)

			srv := newTestServer(t, mockSvc)
			req := &authv1.RegisterRequest{User: amprofiles.UsrToProtoUsr(user)}

			_, err := srv.Register(context.Background(), req)
			assert.Error(t, err)
			st, _ := status.FromError(err)
			assert.Equal(t, tt.wantCode, st.Code())
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestIsAdmin_Success(t *testing.T) {
	mockSvc := new(MockAuthService)
	id := uuid.New()

	mockSvc.On("IsAdmin", mock.Anything, id).Return(true, nil)

	srv := newTestServer(t, mockSvc)
	req := &authv1.IsAdminRequest{UserId: id.String()}

	resp, err := srv.IsAdmin(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, resp.IsAdmin)
	mockSvc.AssertExpectations(t)
}

func TestIsAdmin_InvalidUUID(t *testing.T) {
	mockSvc := new(MockAuthService)
	srv := newTestServer(t, mockSvc)

	req := &authv1.IsAdminRequest{UserId: "invalid-uuid"}

	_, err := srv.IsAdmin(context.Background(), req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestIsAdmin_Errors(t *testing.T) {
	mockSvc := new(MockAuthService)
	id := uuid.New()

	tests := []struct {
		name       string
		serviceErr error
		wantCode   codes.Code
	}{
		{"DeadlineExceeded", serviceerrors.ErrDeadlineExceeded, codes.DeadlineExceeded},
		{"InvalidArgument", serviceerrors.ErrInvalidArgument, codes.InvalidArgument},
		{"NotFound", serviceerrors.ErrNotFound, codes.NotFound},
		{"OtherError", errors.New("other"), codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc.ExpectedCalls = nil
			mockSvc.On("IsAdmin", mock.Anything, id).Return(false, tt.serviceErr)

			srv := newTestServer(t, mockSvc)
			req := &authv1.IsAdminRequest{UserId: id.String()}

			_, err := srv.IsAdmin(context.Background(), req)
			assert.Error(t, err)
			st, _ := status.FromError(err)
			assert.Equal(t, tt.wantCode, st.Code())
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestIsAdmin_ContextDone(t *testing.T) {
	mockSvc := new(MockAuthService)
	srv := newTestServer(t, mockSvc)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &authv1.IsAdminRequest{UserId: uuid.New().String()}

	_, err := srv.IsAdmin(ctx, req)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
}
