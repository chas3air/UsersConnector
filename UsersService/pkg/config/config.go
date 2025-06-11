package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                    string `yaml:"env" env-default:"local"`
	Port                   int    `yaml:"port" env:"PORT" env-default:"8080"`
	MongoDBHost            string `yaml:"mongodb_host" env:"MONGODB_HOST" env-default:"mongo_cont"`
	MongoDBPort            int    `yaml:"mongodb_port" env:"MONGODB_PORT" env-default:"27017"`
	MongoDBDBName          string `yaml:"mongodb_db_name" env:"MONGODB_DB_NAME" env-default:"users"`
	MongoDBUsersCollection string `yaml:"mongodb_users_collection_name" env:"MONGODB_USERS_COLLECTION_NAME" env-default:"users"`
	PostgresHost           string `yaml:"psql_host" env:"PSQL_HOST" env-default:"psql_cont"`
	PostgresPort           int    `yaml:"psql_port" env:"PSQL_PORT" env-default:"5432"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
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
