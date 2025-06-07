package main

import (
	"auth/pkg/config"
	"auth/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	storage := grpcusers.New(log, os.Getenv("USERS_HOST"), os.Getenv("USERS_PORT"))

	// application := app.New(cfg, log, storage)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	storage.Close()

	application.GRPCServer.Stop()
}
