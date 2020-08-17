package main

import (
	"github.com/vsmoraes/open-ledger/cmd/rest-api/application"
	"github.com/vsmoraes/open-ledger/config"
)

func main() {
	app := application.NewApplication()
	app.Start(config.Config().RestAPI.Port)
}
