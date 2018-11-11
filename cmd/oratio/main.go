package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"milobella/oratio/pkg/ability"
	"milobella/oratio/pkg/anima"
	"milobella/oratio/pkg/cerebro"
	"net/http"
	"time"
)

var cerebroCli = cerebro.NewClient("http://0.0.0.0", 9444)
var animaCli = anima.NewClient("http://0.0.0.0", 9333)
var cinemaCli = ability.NewClient("http://0.0.0.0", 10200)

// fun main()
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/talk/text", TextRequest).Methods("POST")
    log.Fatal(http.ListenAndServe(":9100", router))
}

type RequestBody struct {
    Text 	string   `json:"text,omitempty"`
}

type ResponseBody struct {
    Vocal  	string `json:"vocal,omitempty"`
}

func TextRequest(w http.ResponseWriter, r *http.Request) {

	// Read the request
	body, err := ReadRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// Execute the processing flow
	var nlu = cerebroCli.UnderstandText(body.Text)
	var nlg = CallAbility(nlu)
	var vox = animaCli.GenerateSentence(nlg)

	// Build the response
	json.NewEncoder(w).Encode(ResponseBody{Vocal: vox})
}

func ReadRequest(r *http.Request) (req RequestBody, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &req)
	return
}

func CallAbility(nlu cerebro.NLU) anima.NLG {
	// TODO put personal request in anima
	if nlu.Intent == "HELLO"{
		return anima.NLG{Sentence: "Hello"}
	}

	// TODO put time request in clock ability
	if nlu.Intent == "GET_TIME" {
		timeVal := fmt.Sprintf("%d h %d", time.Now().Hour(), time.Now().Minute())
		return anima.NLG{Sentence: "It is {{time}}", Params: map[string]string{"time": timeVal}}
	}

	if nlu.Intent == "LAST_SHOWTIME" {
		return cinemaCli.CallAbility(nlu)
	}

	return anima.NLG{Sentence: "Oups !"}
}
