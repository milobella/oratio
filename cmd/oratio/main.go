package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/auth"
	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/logging"
	"github.com/milobella/oratio/internal/persistence"
	"github.com/milobella/oratio/internal/server"
	"github.com/milobella/oratio/internal/service"
	"github.com/milobella/oratio/internal/tracing"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	logrus.SetReportCaller(true)
}

//TODO simplify the main
// - making some builders to initialize handlers
func main() {
	// Read configuration
	conf := config.Read()

	shutdown, err := tracing.InitGlobalTracer(conf.Tracing)
	if err != nil {
		logrus.WithError(err).Fatal("Could not initialize jaeger tracer.")
	}
	defer shutdown()

	// Initialize clients
	cerebroClient := cerebro.NewClient(conf.Cerebro.Host, conf.Cerebro.Port, conf.Cerebro.UnderstandEndpoint)
	animaClient := anima.NewClient(conf.Anima.Host, conf.Anima.Port, conf.Anima.RestituteEndpoint)

	// Build the ability service
	// TODO(Célian): The DAO being not initialized shouldn't be a reason of not being up. We should have a retry mechanism + health endpoint
	abilityDAO, err := persistence.NewAbilityDAOMongo(conf.AbilitiesDatabase, 3*time.Second)
	if err != nil {
		logrus.WithError(err).Fatalf("Error initializing the Ability DAO.")
	}
	abilityService := service.NewAbilityService(abilityDAO, conf.Abilities, conf.AbilitiesCache.Expiration, conf.AbilitiesCache.CleanupInterval)

	// Build the ability handler
	abilityHandler := &server.AbilityRequestHandler{
		AbilitDAO:      abilityDAO,
		AbilityService: abilityService,
	}

	// Build the text handler
	textHandler := server.TextRequestHandler{
		CerebroClient:  cerebroClient,
		AnimaClient:    animaClient,
		AbilityService: abilityService,
	}

	// Init an echo server
	applicationServer := echo.New()

	// Register middlewares
	tracing.ApplyMiddleware(applicationServer, conf.Tracing)
	logging.ApplyMiddleware(applicationServer)
	auth.ApplyMiddleware(applicationServer, conf.Auth)

	// Register handlers
	apiV1 := applicationServer.Group("/api/v1")
	apiV1.POST("/talk/text", textHandler.HandleTextRequest)
	apiV1.GET("/abilities", abilityHandler.HandleGetAllAbilityRequest)
	apiV1.POST("/abilities", abilityHandler.HandleCreateAbilityRequest)

	// Keep the old route to ensure the compatibility
	// TODO: remove old route after the migration is performed
	applicationServer.POST("/talk/text", textHandler.HandleTextRequest)

	// Run the echo server
	logrus.Fatal(applicationServer.Start(fmt.Sprintf(":%d", conf.Server.Port)))
}
