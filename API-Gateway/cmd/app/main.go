package main

import (
	"api-gateway/internal/app"
	grpcauthserver "api-gateway/internal/storage/grpc/auth"
	grpcusersstorage "api-gateway/internal/storage/grpc/users"
	userscashstorage "api-gateway/internal/storage/redis/users"
	"api-gateway/pkg/config"
	"api-gateway/pkg/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadEnv()

	log := logger.SetupLogger(cfg.Env)

	log.Info("application configured", slog.Any("config", cfg))

	grpcUsersApiConnection := grpcusersstorage.New(log, cfg.GrpcUsersAPIHost, cfg.GrpcUsersAPIPort)
	log.Info("connection to usersService done")
	grpcAuthApiConnection := grpcauthserver.New(log, cfg.GrpcAuthAPIHost, cfg.GrpcAuthAPIPort)
	log.Info("connection to authService done")
	redisConnection := userscashstorage.New(log, cfg.RedisHost, cfg.RedisPort, cfg.ExpirationTime)
	log.Info("connection to redis done")

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
