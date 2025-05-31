package usershandler

import (
	"api-gateway/internal/domain/models"
	serviceerror "api-gateway/internal/service"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IUsersService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type UsersHandler struct {
	log     *slog.Logger
	service IUsersService
}

func New(log *slog.Logger, service IUsersService) *UsersHandler {
	return &UsersHandler{
		log:     log,
		service: service,
	}
}

func (u *UsersHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.GetUsersHandler"
	log := u.log.With(
		"op", op,
	)

	users, err := u.service.GetUsers(r.Context())
	if err != nil {
		log.Error("Error fetching users", sl.Err(err))
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Error("Cannot write users to response", sl.Err(err))
		http.Error(w, "Cannot write users to response", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.GetUserByHandler"
	log := u.log.With(
		"op", op,
	)

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("Id is required")
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	user, err := u.service.GetUserById(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Error("Cannot fetch user by id", sl.Err(err))
		http.Error(w, "Cannot fetch user by id", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error("Cannot write user to response", sl.Err(err))
		http.Error(w, "Cannot write user to response", http.StatusInternalServerError)
		return
	}
}
func (u *UsersHandler) InsertHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.InsertHandler"
	log := u.log.With("op", op)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userForInsert models.User
	if err := json.Unmarshal(body, &userForInsert); err != nil {
		log.Error("Cannot parse body to user", sl.Err(err))
		http.Error(w, "Cannot parse body to user", http.StatusBadRequest)
		return
	}

	insertedUser, err := u.service.Insert(r.Context(), userForInsert)
	if err != nil {
		if errors.Is(err, serviceerror.ErrAlreadyExists) {
			log.Warn("User already exists", sl.Err(err))
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		log.Error("Cannot insert user", sl.Err(err))
		http.Error(w, "Cannot insert user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(insertedUser); err != nil {
		log.Error("Cannot write user to response", sl.Err(err))
		http.Error(w, "Cannot write user to response", http.StatusInternalServerError)
		return
	}
}
func (u *UsersHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.UpdateHandler"
	log := u.log.With(
		"op", op,
	)

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("Id is required")
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Cannot read request body", sl.Err(err))
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userForUpdate models.User
	if err := json.Unmarshal(body, &userForUpdate); err != nil {
		log.Error("Cannot parse body to user", sl.Err(err))
		http.Error(w, "Cannot parse body to user", http.StatusBadRequest)
		return
	}

	updatedUser, err := u.service.Update(r.Context(), id, userForUpdate)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Error("User not found", sl.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Error("Cannot update user", sl.Err(err))
		http.Error(w, "Cannot update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		log.Error("Cannot write user to response", sl.Err(err))
		http.Error(w, "Cannot write user to response", http.StatusInternalServerError)
		return
	}
}

func (u *UsersHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.users.DeleteHandler"
	log := u.log.With(
		"op", op,
	)

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Error("Id is required")
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	deletedUser, err := u.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Error("User not found", sl.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Error("Cannot delete user", sl.Err(err))
		http.Error(w, "Cannot delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(deletedUser); err != nil {
		log.Error("Cannot write user to response", sl.Err(err))
		http.Error(w, "Cannot write user to response", http.StatusInternalServerError)
		return
	}
}
