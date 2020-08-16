package main

import (
	"github.com/vsmoraes/open-ledger/cmd/rest-api/application"
	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

var mf movement.FindableByAccount
var l *ledger.Ledger

func main() {
	app := application.NewApplication()
	app.Start(":8000")
}
