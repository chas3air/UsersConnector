package app

import (
	"api-gateway/internal/domain/models"
	authhandler "api-gateway/internal/handlers/auth"
	usershandler "api-gateway/internal/handlers/users"
	authservice "api-gateway/internal/service/auth"
	userscashservice "api-gateway/internal/service/redis/users"
	usersservice "api-gateway/internal/service/users"
	grpcstorage "api-gateway/internal/storage/grpc/users"
	"api-gateway/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IUsersStorage interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error)
	Insert(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error)
	Delete(ctx context.Context, uid uuid.UUID) (models.User, error)
}

type IAuthServer interface {
	Login(ctx context.Context, login string, password string) (string, string, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

type App struct {
	cfg          *config.Config
	log          *slog.Logger
	psqlStorage  IUsersStorage
	authServer   IAuthServer
	redisStorage userscashservice.UsersCashStorage
}

func New(cfg *config.Config, log *slog.Logger, storage *grpcstorage.GRPCUsersStorage, authServer IAuthServer, redisStorage userscashservice.UsersCashStorage) *App {
	return &App{
		cfg:          cfg,
		log:          log,
		psqlStorage:  storage,
		authServer:   authServer,
		redisStorage: redisStorage,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	r := mux.NewRouter()

	redisService := userscashservice.New(a.log, a.redisStorage)

	usersService := usersservice.New(a.log, a.psqlStorage)
	usersHandler := usershandler.New(a.log, usersService, redisService, a.cfg.MaxRequestsPerUser)
	a.log.Info("usersHandler done")

	authService := authservice.New(a.log, a.psqlStorage, a.authServer)
	authHandler := authhandler.New(a.log, authService, redisService)
	a.log.Info("authHandler done")

	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 OK"))
	})
	r.HandleFunc("/api/v1/login", authHandler.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/register", authHandler.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/refresh", authHandler.RefreshTokenHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/logout", authHandler.LoginHandler).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/users", usersHandler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", usersHandler.InsertHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.UpdateHandler).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.DeleteHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Port), r); err != nil {
		return err
	}

	a.log.Info("application is listened")

	return nil
}
