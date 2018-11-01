package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"milobella/oratio/pkg/anima"
	"milobella/oratio/pkg/cerebro"
	"net/http"
)

var cerebroCli = cerebro.Client{}
var animaCli = anima.Client{}

// fun main()
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/talk/text", TextRequest).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
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
	var nlg = CallSkill(nlu)
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

// TODO: call the skill according to the skill field
func CallSkill(nlu cerebro.NLU) anima.NLG {
	if nlu.Skill == "hello" && nlu.Action == "hello" {
		return anima.NLG{Sentence: "Bonjour"}
	}
	return anima.NLG{Sentence: "Erreur"}
}
