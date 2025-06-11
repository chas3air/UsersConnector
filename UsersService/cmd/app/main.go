package main

import (
	"os"
	"os/signal"
	"syscall"
	"usersservice/internal/app"
	usersstorage "usersservice/internal/storage/mongo/users"
	"usersservice/pkg/config"
	"usersservice/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	storage := usersstorage.New(log, cfg.MongoDBHost, cfg.MongoDBPort, cfg.MongoDBDBName, cfg.MongoDBUsersCollection)

	application := app.New(log, cfg.Port, storage)
	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	storage.Close()
	log.Info("Database connection closed")

	application.GRPCServer.Stop()
	log.Info("application stoped")
}
