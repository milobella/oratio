package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/celian-garcia/gonfig"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gitlab.milobella.com/milobella/oratio/internal/config"
	"gitlab.milobella.com/milobella/oratio/pkg/ability"
	"gitlab.milobella.com/milobella/oratio/pkg/anima"
	"gitlab.milobella.com/milobella/oratio/pkg/cerebro"
)

//TODO: use this init function to initialize variables instead of initialize on top
func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

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
		abilityClient := ability.NewClient(ac.Host, ac.Port)
		for _, intent := range ac.Intents {
			abilityClientsMap[intent] = abilityClient
		}
	}

	// Initialize the server's router
	router := mux.NewRouter()
	router.HandleFunc("/talk/text", textRequest).Methods("POST")

	// Initializing the server
	addr := fmt.Sprintf(":%d", conf.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("Error initializing the server : %s", err)
	}

	// Start the server
	done := make(chan bool)
	go http.Serve(listener, router)
	logrus.Infof("Successfully started the Milobella::Oratio server on port %d !", conf.Server.Port)
	<-done
}

type RequestBody struct {
	Text string `json:"text,omitempty"`
}

type ResponseBody struct {
	Vocal        string      `json:"vocal,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
}

func textRequest(w http.ResponseWriter, r *http.Request) {

	// Read the request
	body, err := readRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// Execute the processing flow
	nlu := cerebroClient.UnderstandText(body.Text)
	nlg, visu, auto_reprompt := callAbility(nlu)
	vocal := animaClient.GenerateSentence(nlg)

	// Build the response
	json.NewEncoder(w).Encode(ResponseBody{Vocal: vocal, Visu: visu, AutoReprompt: auto_reprompt})
}

func readRequest(r *http.Request) (req RequestBody, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &req)
	return
}

// Choose what ability to call according to the intent resolved by cerebro.
func callAbility(nlu cerebro.NLU) (nlg anima.NLG, visu interface{}, auto_reprompt bool) {

	// TODO put personal request in anima
	if nlu.BestIntent == "HELLO" {
		return anima.NLG{Sentence: "Hello"}, nil, false
	}

	// TODO put time request in clock ability
	if nlu.BestIntent == "GET_TIME" {
		timeVal := fmt.Sprintf("%d h %d", time.Now().Hour(), time.Now().Minute())
		return anima.NLG{
			Sentence: "It is {{time}}",
			Params: []anima.NLGParam{{
				Name:  "time",
				Value: timeVal,
				Type:  "time",
			}}}, nil, false
	}

	if client, ok := abilityClientsMap[nlu.BestIntent]; ok {
		return client.CallAbility(nlu)
	}

	return anima.NLG{Sentence: "Oups !"}, nil, false
}
