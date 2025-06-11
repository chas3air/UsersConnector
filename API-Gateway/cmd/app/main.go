package main

import (
	"api-gateway/internal/app"
	grpcstorage "api-gateway/internal/storage/grpc/users"
	userscashstorage "api-gateway/internal/storage/redis/users"
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

	grpcConnection := grpcstorage.New(log, cfg.GrpcUsersAPIHost, cfg.GrpcUsersAPIPort)
	redisConnection := userscashstorage.New(log, cfg.RedisHost, cfg.RedisPort, cfg.ExpirationTime)

	application := app.New(cfg, log, grpcConnection, redisConnection)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	grpcConnection.Close()
	log.Info("grpcConnection closed")

	redisConnection.Close()
	log.Info("redisConnection closed")

	log.Info("application stoped")
}
