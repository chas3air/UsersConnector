package userspsqlstorage_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"usersservice/internal/domain/models"
	storageerror "usersservice/internal/storage"
	userspsqlstorage "usersservice/internal/storage/psql/users"
	"usersservice/pkg/lib/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func newTestStorage(t *testing.T) (*userspsqlstorage.UsersPsqlStorage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	storage := &userspsqlstorage.UsersPsqlStorage{
		Log:       logger.SetupLogger("local"),
		DB:        db,
		TableName: "users",
	}

	cleanup := func() {
		db.Close()
	}

	return storage, mock, cleanup
}

func TestGetUsers(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "login", "password", "role"}).
		AddRow(uuid.New(), "user1", "pass1", "admin").
		AddRow(uuid.New(), "user2", "pass2", "user")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users;")).WillReturnRows(rows)

	users, err := storage.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserById_Success(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	id := uuid.New()
	row := sqlmock.NewRows([]string{"id", "login", "password", "role"}).
		AddRow(id, "user1", "pass1", "admin")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id=$1;")).
		WithArgs(id).
		WillReturnRows(row)

	user, err := storage.GetUserById(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Id != id {
		t.Errorf("expected id %v, got %v", id, user.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id=$1;")).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := storage.GetUserById(context.Background(), id)
	if !errors.Is(err, storageerror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert_Success(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	user := models.User{
		Id:       uuid.New(),
		Login:    "user1",
		Password: "pass1",
		Role:     "admin",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, login, password, role) VALUES ($1, $2, $3, $4);")).
		WithArgs(user.Id, user.Login, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertedUser, err := storage.Insert(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insertedUser.Id != user.Id {
		t.Errorf("expected id %v, got %v", user.Id, insertedUser.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert_AlreadyExists(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	user := models.User{
		Id:       uuid.New(),
		Login:    "user1",
		Password: "pass1",
		Role:     "admin",
	}

	pqErr := &pq.Error{Code: "23505"}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, login, password, role) VALUES ($1, $2, $3, $4);")).
		WithArgs(user.Id, user.Login, user.Password, user.Role).
		WillReturnError(pqErr)

	_, err := storage.Insert(context.Background(), user)
	if !errors.Is(err, storageerror.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	user := models.User{
		Id:       uuid.New(),
		Login:    "user1",
		Password: "pass1",
		Role:     "admin",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET login=$1, password=$2, role=$3 WHERE id=$4;")).
		WithArgs(user.Login, user.Password, user.Role, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedUser, err := storage.Update(context.Background(), user.Id, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedUser.Id != user.Id {
		t.Errorf("expected id %v, got %v", user.Id, updatedUser.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	user := models.User{
		Id:       uuid.New(),
		Login:    "user1",
		Password: "pass1",
		Role:     "admin",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET login=$1, password=$2, role=$3 WHERE id=$4;")).
		WithArgs(user.Login, user.Password, user.Role, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 0))

	_, err := storage.Update(context.Background(), user.Id, user)
	if !errors.Is(err, storageerror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete_Success(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	id := uuid.New()
	user := models.User{
		Id:       id,
		Login:    "user1",
		Password: "pass1",
		Role:     "admin",
	}

	row := sqlmock.NewRows([]string{"id", "login", "password", "role"}).
		AddRow(user.Id, user.Login, user.Password, user.Role)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id=$1;")).
		WithArgs(id).
		WillReturnRows(row)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1;")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deletedUser, err := storage.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedUser.Id != id {
		t.Errorf("expected id %v, got %v", id, deletedUser.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	storage, mock, cleanup := newTestStorage(t)
	defer cleanup()

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id=$1;")).
		WithArgs(id).
		WillReturnError(storageerror.ErrNotFound)

	_, err := storage.Delete(context.Background(), id)
	if !errors.Is(err, storageerror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
