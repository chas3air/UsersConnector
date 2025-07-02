package authhandler

import (
	"api-gateway/internal/domain/models"
	serviceerror "api-gateway/internal/service"
	"api-gateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const userRoleTitle = "user"

// const adminRoleTitle = "admin"

// var roleRating = map[string]int{
// 	userRoleTitle: 4,
// 	adminRoleTitle: 8,
// }

type IAuthService interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

type IUserCashService interface {
	Get(context.Context, uuid.UUID) (models.User, error)
	Set(context.Context, models.User) error
	Del(context.Context, uuid.UUID) error
}

type AuthHandler struct {
	log          *slog.Logger
	service      IAuthService
	redisService IUserCashService
}

func New(log *slog.Logger, service IAuthService, redisService IUserCashService) *AuthHandler {
	return &AuthHandler{
		log:          log,
		service:      service,
		redisService: redisService,
	}
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Login"
	log := a.log.With(
		"op", op,
	)

	var loginStruct = struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&loginStruct); err != nil {
		log.Error("Cannot parse request body to obj", sl.Err(err))
		http.Error(w, "Cannot parse request body to obj", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := a.service.Login(r.Context(), loginStruct.Login, loginStruct.Password)
	if err != nil {
		if errors.Is(err, serviceerror.ErrNotFound) {
			log.Warn("User not found", sl.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Error("Cannot login", sl.Err(err))
		http.Error(w, "Cannot login", http.StatusInternalServerError)
		return
	}

	tokenResponse := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		log.Error("Cannot write token to response", sl.Err(err))
		http.Error(w, "Cannot write token to response", http.StatusInternalServerError)
		return
	}
}

func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Register"
	log := a.log.With(
		"op", op,
	)

	var userForRegister models.User
	if err := json.NewDecoder(r.Body).Decode(&userForRegister); err != nil {
		log.Error("Cannot read requesy body", sl.Err(err))
		http.Error(w, "Cannot read requesy body", http.StatusBadRequest)
		return
	}
	userForRegister.Role = userRoleTitle

	registeredUser, err := a.service.Register(r.Context(), userForRegister)
	if err != nil {
		if errors.Is(err, serviceerror.ErrAlreadyExists) {
			log.Warn("User already registered", sl.Err(err))
			http.Error(w, "User already registered", http.StatusConflict)
			return
		}

		log.Error("Cannot register", sl.Err(err))
		http.Error(w, "Cannot register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(registeredUser.Id); err != nil {
		log.Error("Cannot write id to response", sl.Err(err))
		http.Error(w, "Cannot write id to response", http.StatusInternalServerError)
		return
	}
}

func (AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {}
func (AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request)       {}
