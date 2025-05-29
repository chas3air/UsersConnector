package usersstorage

import (
	"context"
	"database/sql"
	"log/slog"
	"usersservice/internal/domain/models"
	"usersservice/pkg/lib/logger/sl"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UsersStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

func New(log *slog.Logger, connStr string) *UsersStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Error connecting to database", sl.Err(err))
		panic(err)
	}

	/*
		wd, _ := os.GetWd()
			migrationPath := filepath.Join(wd, "app", "migrations")
		if err := applyMigrations(db, migrationPath); err != nil {
			panic(err)
		}
	*/

	return &UsersStorage{
		log: log,
		DB:  db,
	}
}

func (u *UsersStorage) Close() {
	if err := u.DB.Close(); err != nil {
		panic(err)
	}
}

// GetUsers implements app.IUsersStorage.
func (u *UsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	panic("unimplemented")
}

// GetUserById implements app.IUsersStorage.
func (u *UsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// Insert implements app.IUsersStorage.
func (u *UsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Update implements app.IUsersStorage.
func (u *UsersStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	panic("unimplemented")
}

// Delete implements app.IUsersStorage.
func (u *UsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	panic("unimplemented")
}
