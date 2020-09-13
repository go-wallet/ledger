package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/vsmoraes/open-ledger/cmd/rest-api/application"
	"github.com/vsmoraes/open-ledger/internal/config"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	app := application.NewApplication()
	app.Start(config.Config().RestAPI.Port)
}
