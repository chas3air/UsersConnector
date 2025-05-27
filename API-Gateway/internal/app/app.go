package app

import (
	authhandler "api-gateway/internal/handlers/auth"
	usershandler "api-gateway/internal/handlers/users"
	authservice "api-gateway/internal/service/auth"
	usersservice "api-gateway/internal/service/users"
	usersstorage "api-gateway/internal/storage/users"
	"api-gateway/pkg/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	r := mux.NewRouter()

	usersStorage := usersstorage.New(a.log, "users", 1234)
	usersService := usersservice.New(a.log, usersStorage)
	usersHandler := usershandler.New(a.log, usersService)

	authService := authservice.New(a.log, usersStorage)
	authHandler := authhandler.New(a.log, authService)

	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("200 OK"))
	})
	r.HandleFunc("/api/v1/login", authHandler.LoginHandler)
	r.HandleFunc("/api/v1/register", authHandler.RegisterHandler)
	r.HandleFunc("/api/v1/refresh", authHandler.RefreshTokenHandler)
	r.HandleFunc("/api/v1/logout", authHandler.LoginHandler)

	r.HandleFunc("/api/v1/users", usersHandler.GetUsersHandler)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.GetUserByIdHandler)
	r.HandleFunc("/api/v1/users", usersHandler.InsertHandler)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.UpdateHandler)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.DeleteHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Port), r); err != nil {
		return err
	}

	return nil
}
