package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	Port           int    `yaml:"port" env:"PORT" env-default:"8080"`
	ExpirationTime int    `yaml:"expiration_time" env:"EXPIRATION_TIME" env-default:"10"`

	GrpcUsersAPIHost string `yaml:"grpc_users_api_host" env:"GRPC_USERS_API_HOST" env-default:"usersservice"`
	GrpcUsersAPIPort int    `yaml:"grpc_users_api_port" env:"GRPC_USERS_API_PORT" env-default:"50051"`

	GrpcAuthAPIHost string `yaml:"grpc_auth_api_host" env:"GRPC_AUTH_API_HOST" env-default:"auth"`
	GrpcAuthAPIPort int    `yaml:"grpc_auth_api_port" env:"GRPC_AUTH_API_PORT" env-default:"50051"`

	RedisHost string `yaml:"redis_host" env:"REDIS_HOST" env-default:"redis"`
	RedisPort int    `yaml:"redis_port" env:"REDIS_PORT" env-default:"6379"`

	MaxRequestsPerUser int `yaml:"max_requests_per_user" env:"MAX_REQUESTS_PER_USER" env-default:"100"`
}

func MustLoadYaml() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println(os.Getwd())
		log.Println("Error loading .env file")
		panic(err)
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("cannot read config from environment: " + err.Error())
	}

	return &cfg
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	// --config=./config/local.yaml
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
