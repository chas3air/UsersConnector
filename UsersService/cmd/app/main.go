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

	storage := usersstorage.New(log, "mongodb://mongo_cont:27017")

	application := app.New(log, cfg.Grpc.Port, storage)
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
