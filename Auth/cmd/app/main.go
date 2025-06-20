package main

import (
	"auth/internal/app"
	grpcusers "auth/internal/storage/grpc/users"
	"auth/pkg/config"
	"auth/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	usersConnection := grpcusers.New(log, cfg.GrpcUsersAPIHost, cfg.GrpcUsersAPIPort)

	application := app.New(log, cfg.Grpc.Port, usersConnection)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	usersConnection.Close()

	application.GRPCServer.Stop()
}
