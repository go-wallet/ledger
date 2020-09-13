package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

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

type redisInstance struct {
	Host     string
	Password string
}

type redisConfig struct {
	Instances []redisInstance
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
			_, b, _, _ := runtime.Caller(0)
			dir := path.Join(path.Dir(b))

			godotenv.Load(fmt.Sprintf("%s/../../.env", dir))
		}

		port := fmt.Sprintf(":%s", os.Getenv("APP_API_PORT"))
		mongoURI := os.Getenv("APP_MONGODB_URI")
		redisHosts := os.Getenv("APP_REDIS_HOSTS")
		redisPasswords := os.Getenv("APP_REDIS_PASSWORDS")

		configInstance = &config{
			RestAPI: &restAppConfig{
				Port: port,
			},

			MongoDB: &mongoConfig{
				URI:                 mongoURI,
				Database:            "open-ledger",
				MovementsCollection: "movements",
			},

			Redis: &redisConfig{
				Instances: parseRegisConfig(redisHosts, redisPasswords),
			},
		}
	}

	return configInstance
}

func parseRegisConfig(hosts, passwords string) []redisInstance {
	hs := strings.Split(hosts, ",")
	ps := strings.Split(passwords, ",")

	instances := make([]redisInstance, 0)
	for i, host := range hs {
		instances = append(instances, redisInstance{
			Host:     host,
			Password: ps[i],
		})
	}

	return instances
}
