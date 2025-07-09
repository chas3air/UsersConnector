package authservice_test

import (
	"context"
	"errors"
	"testing"

	"auth/internal/domain/models"
	serviceerrors "auth/internal/service"
	authservice "auth/internal/service/auth"
	storageerrors "auth/internal/storage"
	"auth/pkg/lib/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock IUsersStorage ---

type MockUsersStorage struct {
	mock.Mock
}

func (m *MockUsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

// --- Tests ---

func newTestService(storage *MockUsersStorage) *authservice.AuthService {
	logger := logger.SetupLogger("local")
	return authservice.New(logger, storage)
}

func TestLogin_Success(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	user := models.User{Login: "testuser", Password: "pass123"}
	mockStorage.On("GetUsers", mock.Anything).Return([]models.User{user}, nil)

	svc := newTestService(mockStorage)

	access, refresh, err := svc.Login(context.Background(), "testuser", "pass123")
	assert.NoError(t, err)
	assert.Equal(t, "access-token", access)
	assert.Equal(t, "refresh-token", refresh)
	mockStorage.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	mockStorage.On("GetUsers", mock.Anything).Return([]models.User{}, nil)

	svc := newTestService(mockStorage)

	_, _, err := svc.Login(context.Background(), "no-user", "pass")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user doesn't exists")
	mockStorage.AssertExpectations(t)
}

func TestLogin_GetUsersError(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	mockStorage.On("GetUsers", mock.Anything).Return([]models.User(nil), errors.New("db error"))

	svc := newTestService(mockStorage)

	_, _, err := svc.Login(context.Background(), "any", "any")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockStorage.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	newUser := models.User{Login: "newuser", Password: "pass123"}

	// В хранилище нет пользователей с таким логином и паролем
	mockStorage.On("GetUsers", mock.Anything).Return([]models.User{}, nil)
	mockStorage.On("Insert", mock.Anything, newUser).Return(newUser, nil)

	svc := newTestService(mockStorage)

	inserted, err := svc.Register(context.Background(), newUser)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Login, inserted.Login)
	mockStorage.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	existingUser := models.User{Login: "existuser", Password: "pass123"}

	mockStorage.On("GetUsers", mock.Anything).Return([]models.User{existingUser}, nil)

	svc := newTestService(mockStorage)

	_, err := svc.Register(context.Background(), existingUser)
	assert.ErrorIs(t, err, serviceerrors.ErrAlreadyExists)
	mockStorage.AssertExpectations(t)
}

func TestRegister_InsertError(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	newUser := models.User{Login: "newuser", Password: "pass123"}

	mockStorage.On("GetUsers", mock.Anything).Return([]models.User{}, nil)
	mockStorage.On("Insert", mock.Anything, newUser).Return(models.User{}, errors.New("insert error"))

	svc := newTestService(mockStorage)

	_, err := svc.Register(context.Background(), newUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert error")
	mockStorage.AssertExpectations(t)
}

func TestIsAdmin_Success(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	id := uuid.New()
	adminUser := models.User{Id: id, Role: "admin"}

	mockStorage.On("GetUserById", mock.Anything, id).Return(adminUser, nil)

	svc := newTestService(mockStorage)

	isAdmin, err := svc.IsAdmin(context.Background(), id)
	assert.NoError(t, err)
	assert.True(t, isAdmin)
	mockStorage.AssertExpectations(t)
}

func TestIsAdmin_NotAdmin(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	id := uuid.New()
	user := models.User{Id: id, Role: "user"}

	mockStorage.On("GetUserById", mock.Anything, id).Return(user, nil)

	svc := newTestService(mockStorage)

	isAdmin, err := svc.IsAdmin(context.Background(), id)
	assert.NoError(t, err)
	assert.False(t, isAdmin)
	mockStorage.AssertExpectations(t)
}

func TestIsAdmin_GetUserByIdErrors(t *testing.T) {
	mockStorage := new(MockUsersStorage)
	id := uuid.New()

	tests := []struct {
		name        string
		storageErr  error
		expectedErr error
		expectedLog string
	}{
		{"DeadlineExceeded", storageerrors.ErrDeadlineExceeded, serviceerrors.ErrDeadlineExceeded, "Deadline exceeded"},
		{"InvalidArgument", storageerrors.ErrInvalidArgument, serviceerrors.ErrInvalidArgument, "Invalid argument"},
		{"NotFound", storageerrors.ErrNotFound, serviceerrors.ErrNotFound, "User not found"},
		{"OtherError", errors.New("other error"), errors.New("other error"), "Cannot retrieve user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.ExpectedCalls = nil
			mockStorage.On("GetUserById", mock.Anything, id).Return(models.User{}, tt.storageErr)

			svc := newTestService(mockStorage)

			isAdmin, err := svc.IsAdmin(context.Background(), id)
			assert.False(t, isAdmin)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr.Error())
			mockStorage.AssertExpectations(t)
		})
	}
}
