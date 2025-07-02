package main

import (
	"auth/internal/app"
	grpcusers "auth/internal/storage/grpc/users"
	"auth/pkg/config"
	"auth/pkg/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadEnv()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application config", slog.Any("config", cfg))

	usersConnection := grpcusers.New(log, cfg.GrpcUsersAPIHost, cfg.GrpcUsersAPIPort)

	log.Info("connection configured")

	application := app.New(log, cfg.Port, usersConnection)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	usersConnection.Close()
	log.Info("connection closed")

	application.GRPCServer.Stop()
	log.Info("application stoped")
}
