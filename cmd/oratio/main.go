package main

import (
	"github.com/gorilla/mux"
	"github.com/milobella/oratio/internal/auth"
	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/logging"
	"github.com/milobella/oratio/internal/models"
	"github.com/milobella/oratio/internal/server"
	"github.com/milobella/oratio/pkg/ability"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

//TODO: use this init function to initialize variables instead of initialize on top
func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)
}

var cerebroClient *cerebro.Client
var animaClient *anima.Client
var abilityClientsMap map[string]*ability.Client

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
	cerebroClient = cerebro.NewClient(conf.Cerebro.Host, conf.Cerebro.Port)
	animaClient = anima.NewClient(conf.Anima.Host, conf.Anima.Port)
	abilityClientsMap = make(map[string]*ability.Client)
	// TODO: Abilities should be a dynamic data
	for _, ac := range conf.Abilities {
		abilityClient := ability.NewClient(ac.Host, ac.Port, ac.Name)
		for _, intent := range ac.Intents {
			abilityClientsMap[intent] = abilityClient
		}
	}

	// Build the ability handler
	abilityDAO, err := models.NewAbilityDAOMongo(conf.AbilitiesDatabase.MongoUrl, "oratio", 10 * time.Second)
	if err != nil {
		logrus.WithError(err).Fatalf("Error initializing the Ability DAO.")
	}
	abilityHandler := &server.AbilityRequestHandler{
		AbilitDAO: abilityDAO,
	}

	// Build the text handler
	// TODO: for the moment, abilities are only taken from the configuration but it should take it from the database
	abilityService := &server.AbilityService{Clients: abilityClientsMap}
	textHandler := server.TextRequestHandler{
		CerebroClient:  cerebroClient,
		AnimaClient:    animaClient,
		AbilityService: abilityService,
	}

	// Initialize the server's router
	router := mux.NewRouter()

	// Register the logging requests middleware
	logMiddleware := logging.InitializeLoggingMiddleware()
	router.Use(logMiddleware.Handle)

	// Register the JWT authentication middleware
	if len(conf.AppSecret) > 0 {
		jwtMiddleware := auth.InitializeJWTMiddleware(conf.AppSecret)
		router.Use(jwtMiddleware.Handler)
	}

	router.HandleFunc("/talk/text", textHandler.HandleTextRequest).Methods("POST")
	router.HandleFunc("/abilities", abilityHandler.HandleAbilityRequest).Methods("POST", "GET")

	srv := server.Server{Router: router, Port: conf.Server.Port}
	srv.Run()
}
