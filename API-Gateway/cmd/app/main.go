package main

import (
	"api-gateway/internal/app"
	grpcstorage "api-gateway/internal/storage/grpc/users"
	"api-gateway/pkg/config"
	"api-gateway/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application configured")

	connection := grpcstorage.New(log, os.Getenv("HOST"), os.Getenv("PORT"))

	application := app.New(cfg, log, connection)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	connection.Close()
	log.Info("connection closed")

	log.Info("application stoped")
}
