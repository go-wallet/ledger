package application

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vsmoraes/open-ledger/internal/factory"
	"github.com/vsmoraes/open-ledger/internal/storage"
	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type Application struct {
	e *echo.Echo

	mc *mongo.Client

	mongoClient *storage.MongoClient
	locker      *account.Locker

	mf movement.FindableByAccount
	l  *ledger.Ledger
}

func NewApplication() *Application {
	mongoClient, mc := factory.NewDBRepository()
	locker := factory.NewLocker()

	return &Application{
		e:           echo.New(),
		mc:          mc,
		mongoClient: mongoClient,
		locker:      locker,
		mf:          mongoClient,
		l: ledger.New(
			locker,
			mongoClient,
			mongoClient,
		),
	}
}

func (app *Application) Start(port string) {
	defer app.stop()

	app.e.HideBanner = true
	app.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is up and running")
	})

	rg := app.e.Group("/ledger/movements")
	rg.POST("", createMovementController(app.l))
	rg.GET("", findMomentsController(app.mf))

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Starting HTTP server")
	app.e.Logger.Fatal(app.e.Start(port))
}

func (app *Application) stop() {
	log.Info("Disconnecting...")
	app.mc.Disconnect(context.Background())
}
