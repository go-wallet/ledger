package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type restAppConfig struct {
	Port string
}

type mongoConfig struct {
	URI                 string
	Database            string
	MovementsCollection string
}

type redisConfig struct {
	Host     string
	Password string
}

type config struct {
	RestAPI *restAppConfig
	MongoDB *mongoConfig
	Redis   *redisConfig
}

var configInstance *config = nil

func Config() *config {
	if configInstance == nil {
		env := os.Getenv("APP_ENV")
		if env != "production" {
			godotenv.Load(".env")
		}

		configInstance = &config{
			RestAPI: &restAppConfig{
				Port: fmt.Sprintf(":%s", os.Getenv("APP_API_PORT")),
			},

			MongoDB: &mongoConfig{
				URI:                 os.Getenv("APP_MONGODB_URI"),
				Database:            "open-ledger",
				MovementsCollection: "movements",
			},

			Redis: &redisConfig{
				Host:     os.Getenv("APP_REDIS_HOST"),
				Password: os.Getenv("APP_REDIS_PASSWORD"),
			},
		}
	}

	return configInstance
}
