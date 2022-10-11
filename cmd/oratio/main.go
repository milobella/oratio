package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/auth"
	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/handler"
	"github.com/milobella/oratio/internal/logging"
	"github.com/milobella/oratio/internal/tracing"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)

	logrus.SetReportCaller(true)
}

func main() {
	// Read configuration
	conf := config.Read()

	shutdown, err := tracing.InitGlobalTracer(conf.Tracing)
	if err != nil {
		logrus.WithError(err).Fatal("Could not initialize jaeger tracer.")
	}
	defer shutdown()

	// Init an echo server
	server := echo.New()

	// Create and register middlewares
	tracing.ApplyMiddleware(server, conf.Tracing)
	logging.ApplyMiddleware(server)
	auth.ApplyMiddleware(server, conf.Auth)

	// Create and register handlers
	handlers := handler.New(conf)
	apiV1 := server.Group("/api/v1")
	apiV1.POST("/talk/text", handlers.Text)
	apiV1.GET("/abilities", handlers.GetAbilities)
	apiV1.POST("/abilities", handlers.CreateAbility)

	// Run the echo server
	logrus.Fatal(server.Start(fmt.Sprintf(":%d", conf.Server.Port)))
}
