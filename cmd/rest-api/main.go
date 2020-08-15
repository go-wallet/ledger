package main

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

var mf movement.FindableByAccount
var l *ledger.Ledger

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is up and running")
	})

	rg := e.Group("/ledger/movements")
	rg.POST("", createMovementController(l))
	rg.GET("", findMomentsController(mf))

	e.Logger.Fatal(e.Start(":8000"))
}
