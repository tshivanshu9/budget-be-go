package main

import (
	"fmt"
	"log/slog"

	// "net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tshivanshu9/budget-be/cmd/api/handlers"
	"github.com/tshivanshu9/budget-be/cmd/api/middlewares"
	"github.com/tshivanshu9/budget-be/common"
	"github.com/tshivanshu9/budget-be/internal/mailer"
)

type Application struct {
	logger        *slog.Logger
	server        *echo.Echo
	handler       *handlers.Handler
	appMiddleware middlewares.AppMiddleware
}

func main() {
	e := echo.New()
	db, err := common.NewMysql()

	if err != nil {
		e.Logger.Error("Error with db connection", "error", err)
		return
	}

	err = godotenv.Load()

	if err != nil {
		e.Logger.Error("Error loading .env file", "error", err)
		return
	}
	e.Use(middleware.RequestLogger())
	e.Use(middlewares.CustomMiddleware)

	// e.GET("/", func(c *echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	appMailer := mailer.NewMailer()

	h := &handlers.Handler{
		DB:     db,
		Mailer: appMailer,
	}

	appMiddleware := middlewares.AppMiddleware{
		DB: db,
	}

	app := Application{
		logger:        slog.Default(),
		server:        e,
		handler:       h,
		appMiddleware: appMiddleware,
	}

	// e.GET("/", func(c echo.Context) error {
	// 	return h.Healthcheck(c)
	// })

	// e.GET("/", handlers.Healthcheck(h))
	app.routes()

	fmt.Println(app)

	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)

	if err := e.Start(appAddress); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
