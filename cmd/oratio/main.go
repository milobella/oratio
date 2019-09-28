package main

import (
	"encoding/json"
	"github.com/celian-garcia/gonfig"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"milobella.com/gitlab/milobella/oratio/internal/config"
	"milobella.com/gitlab/milobella/oratio/internal/logging"
	"milobella.com/gitlab/milobella/oratio/internal/server"
	"milobella.com/gitlab/milobella/oratio/pkg/ability"
	"milobella.com/gitlab/milobella/oratio/pkg/anima"
	"milobella.com/gitlab/milobella/oratio/pkg/cerebro"
	"os"
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

// Configuration : object architectured by gonfig lib to store the configuration when we use
type Configuration struct {
	Server     config.ServerConfiguration
	Cerebro    config.CerebroConfiguration
	Anima      config.AnimaConfiguration
	Abilities  []config.AbilityConfiguration
	ConfigFile string `short:"c"`
}

// fun String() : Serialization function of Configuration (for logging)
func (c Configuration) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Configuration serialization error %s", err)
	}
	return string(b)
}

var conf *Configuration

var cerebroClient *cerebro.Client
var animaClient *anima.Client
var abilityClientsMap map[string]*ability.Client

// fun main()
func main() {

	conf = &Configuration{}

	// Load the configuration from file or parameter or env
	err := gonfig.Load(conf, gonfig.Conf{
		ConfigFileVariable: "configfile", // enables passing --configfile myfile.conf

		FileDefaultFilename: "config/oratio.toml",
		FileDecoder:         gonfig.DecoderTOML,

		EnvPrefix: "ORATIO_",
	})

	if err != nil {
		logrus.Fatalf("Error reading config : %s", err)
	} else {
		logrus.Infof("Successfully readen configuration file : %s", conf.ConfigFile)
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

	// Build the handler
	abilityService := &server.AbilityService{Clients: abilityClientsMap}
	textHandler := server.TextRequestHandler{
		CerebroClient: cerebroClient,
		AnimaClient: animaClient,
		AbilityService: abilityService,
	}

	// Initialize the server's router
	router := mux.NewRouter()

	logMiddleware := logging.InitializeLoggingMiddleware()
	router.Use(logMiddleware.Handle)
	router.HandleFunc("/talk/text", textHandler.HandleTextRequest).Methods("POST")

	srv := server.Server{Router: router, Port: conf.Server.Port}
	srv.Run()
}
