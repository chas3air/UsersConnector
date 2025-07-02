package main

import (
	"api-gateway/internal/app"
	"api-gateway/internal/storage/grpc/auth"
	grpcusersstorage "api-gateway/internal/storage/grpc/users"
	userscashstorage "api-gateway/internal/storage/redis/users"
	"api-gateway/pkg/config"
	"api-gateway/pkg/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadEnv()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application configured")

	grpcUsersApiConnection := grpcusersstorage.New(log, cfg.GrpcUsersAPIHost, cfg.GrpcUsersAPIPort)
	grpcAuthApiConnection := grpcauthserver.New(log, cfg.GrpcAuthAPIHost, cfg.GrpcAuthAPIPort)
	redisConnection := userscashstorage.New(log, cfg.RedisHost, cfg.RedisPort, cfg.ExpirationTime)

	application := app.New(cfg, log, grpcUsersApiConnection, grpcAuthApiConnection, redisConnection)

	go func() {
		application.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	grpcUsersApiConnection.Close()
	log.Info("grpcUsersApiConnection closed")

	redisConnection.Close()
	log.Info("redisConnection closed")

	log.Info("application stoped")
}
