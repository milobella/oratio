package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/milobella/oratio/internal"
	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/logging"
	"github.com/milobella/oratio/internal/models"
	"github.com/milobella/oratio/internal/server"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)
}

//TODO simplify the main
// - making some builders to initialize handlers
// - reading configuration in one line (error handling inside the ReadConfiguration function
func main() {
	// Read configuration
	conf, err := config.ReadConfiguration()
	if err != nil { // Handle errors reading the config file
		logrus.WithError(err).Fatalf("Error reading config.")
	} else {
		logrus.Infof("Successfully readen configuration file")
		logrus.Debugf("-> %+v", conf)
	}

	// Initialize clients
	cerebroClient := cerebro.NewClient(conf.Cerebro.Host, conf.Cerebro.Port, conf.Cerebro.UnderstandEndpoint)
	animaClient := anima.NewClient(conf.Anima.Host, conf.Anima.Port, conf.Anima.RestituteEndpoint)

	// Build the ability service
	abilityDAO, err := models.NewAbilityDAOMongo(conf.AbilitiesDatabase.MongoUrl, "oratio", 3*time.Second)
	if err != nil {
		logrus.WithError(err).Fatalf("Error initializing the Ability DAO.")
	}
	abilityService := internal.NewAbilityService(abilityDAO, conf.Abilities, conf.AbilitiesCache.Expiration, conf.AbilitiesCache.CleanupInterval)

	// Build the ability handler
	abilityHandler := &server.AbilityRequestHandler{
		AbilitDAO: abilityDAO,
		AbilityService: abilityService,
	}

	// Build the text handler
	textHandler := server.TextRequestHandler{
		CerebroClient:  cerebroClient,
		AnimaClient:    animaClient,
		AbilityService: abilityService,
	}

	// Initialize an echo server
	e := echo.New()

	// Register middleware
	e.Use(logging.InitializeLoggingMiddleware().Handle)
	if len(conf.AppSecret) > 0 {
		// TODO: use custom claim to retrieve scopes and other user info (https://echo.labstack.com/cookbook/jwt)
		e.Use(middleware.JWT([]byte(conf.AppSecret)))
	}

	// Register handlers
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/talk/text", textHandler.HandleTextRequest)
	apiV1.GET("/abilities", abilityHandler.HandleGetAllAbilityRequest)
	apiV1.POST("/abilities", abilityHandler.HandleCreateAbilityRequest)

	// Keep the old route to ensure the compatibility
	// TODO: remove old route after the migration is performed
	e.POST("/talk/text", textHandler.HandleTextRequest)

	// Run the echo server
	logrus.Fatal(e.Start(fmt.Sprintf(":%d", conf.Server.Port)))
}
