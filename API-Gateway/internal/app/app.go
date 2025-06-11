package app

import (
	authhandler "api-gateway/internal/handlers/auth"
	usershandler "api-gateway/internal/handlers/users"
	authservice "api-gateway/internal/service/auth"
	userscashservice "api-gateway/internal/service/redis/users"
	usersservice "api-gateway/internal/service/users"
	grpcstorage "api-gateway/internal/storage/grpc/users"
	"api-gateway/pkg/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	cfg          *config.Config
	log          *slog.Logger
	psqlStorage  *grpcstorage.GRPCUsersStorage
	redisStorage userscashservice.UsersCashStorage
}

func New(cfg *config.Config, log *slog.Logger, storage *grpcstorage.GRPCUsersStorage, redisStorage userscashservice.UsersCashStorage) *App {
	return &App{
		cfg:          cfg,
		log:          log,
		psqlStorage:  storage,
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

	authService := authservice.New(a.log, a.psqlStorage)
	authHandler := authhandler.New(a.log, authService, redisService)

	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 OK"))
	})
	r.HandleFunc("/api/v1/login", authHandler.LoginHandler)
	r.HandleFunc("/api/v1/register", authHandler.RegisterHandler)
	r.HandleFunc("/api/v1/refresh", authHandler.RefreshTokenHandler)
	r.HandleFunc("/api/v1/logout", authHandler.LoginHandler)

	r.HandleFunc("/api/v1/users", usersHandler.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", usersHandler.InsertHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.UpdateHandler).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.DeleteHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Port), r); err != nil {
		return err
	}

	return nil
}
