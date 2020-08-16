package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vsmoraes/open-ledger/factory"
	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/movement"
	"github.com/vsmoraes/open-ledger/storage"
)

type Application struct {
	e *echo.Echo

	mc *mongo.Client
	rc *redis.Client

	mongoClient *storage.MongoClient
	redisClient *storage.RedisClient

	mf movement.FindableByAccount
	l  *ledger.Ledger
}

func NewApplication() *Application {
	mongoClient, mc := factory.NewDBRepository()
	redisClient, rc := factory.NewLocker()

	return &Application{
		e:           echo.New(),
		mc:          mc,
		rc:          rc,
		mongoClient: mongoClient,
		redisClient: redisClient,
		mf:          mongoClient,
		l: ledger.New(
			redisClient,
			mongoClient,
			mongoClient,
		),
	}
}

func (app *Application) Start(port string) {
	defer app.stop()

	app.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is up and running")
	})

	rg := app.e.Group("/ledger/movements")
	rg.POST("", createMovementController(app.l))
	rg.GET("", findMomentsController(app.mf))

	app.e.Logger.Fatal(app.e.Start(port))
}

func (app *Application) stop() {
	fmt.Println("disconnecting...")
	app.mc.Disconnect(context.Background())
}
