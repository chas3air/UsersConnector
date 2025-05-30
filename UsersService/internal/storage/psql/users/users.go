package usersstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"usersservice/internal/domain/models"
	storageerror "usersservice/internal/storage"
	"usersservice/pkg/lib/logger/sl"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type UsersStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

const UsersTableName = "users"

func New(log *slog.Logger, connStr string) *UsersStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Error connecting to database", sl.Err(err))
		panic(err)
	}

	wd, _ := os.Getwd()
	migrationPath := filepath.Join(wd, "app", "migrations")
	if err := applyMigrations(db, migrationPath); err != nil {
		panic(err)
	}

	return &UsersStorage{
		log: log,
		DB:  db,
	}
}

func applyMigrations(db *sql.DB, migrationsPath string) error {
	return goose.Up(db, migrationsPath)
}

func (u *UsersStorage) Close() {
	if err := u.DB.Close(); err != nil {
		panic(err)
	}
}

// GetUsers implements IUsersStorage.
func (u *UsersStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.psql.users.GetUsers"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	rows, err := u.DB.QueryContext(ctx, `
		SELECT * FROM `+UsersTableName+`;
	`)
	if err != nil {
		log.Error("Error retrieving all users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]models.User, 0, 5)
	var user models.User

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.Role)
		if err != nil {
			log.Warn("Error scanning row", sl.Err(err))
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserById implements IUsersStorage.
func (u *UsersStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.users.GetUserById"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	err := u.DB.QueryRowContext(ctx, `
		SELECT * FROM `+UsersTableName+`
		WHERE id=$1;
	`, uid).Scan(&user.Id, &user.Login, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("User with current id not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
		}

		log.Error("Error scaning row", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Insert implements IUsersStorage.
func (u *UsersStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.psql.users.Insert"
	log := u.log.With(
		"op", op,
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := u.DB.ExecContext(ctx, `
		INSERT INTO `+UsersTableName+` (id, login, password, role)
		VALUES ($1, $2, $3, $4);
	`, user.Id, user.Login, user.Password, user.Role)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Error("User with current id already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrAlreadyExists)
		}

		log.Error("Error inserting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Update implements IUsersStorage.
func (u *UsersStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	const op = "storage.psql.users.Update"
	log := u.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	result, err := u.DB.ExecContext(ctx, `
		UPDATE `+UsersTableName+`
		SET login=$1, password=$2, role=$3
		WHERE id=$4;
	`, user.Login, user.Password, user.Role, user.Id)
	if err != nil {
		log.Error("Error updating user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Error get rows affected", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Error("Zero rows affected")
		return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
	}

	return user, nil
}

// Delete implements IUsersStorage.
func (u *UsersStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.users.Delete"
	log := u.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := u.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storageerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storageerror.ErrNotFound)
		}

		log.Error("Error getting user before deliting", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.DB.ExecContext(ctx, `
		DELETE FROM `+UsersTableName+` 
		WHERE id = $1;
	`, uid)
	if err != nil {
		log.Error("Error deleting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
