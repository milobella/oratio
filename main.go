package main

import (
	"./cerebro"
	"encoding/json"
	"github.com/gorilla/mux"
	"gondora/command-center/anima"
	"io/ioutil"
	"log"
	"net/http"
)

var c = cerebro.CerebroClient{}

// fun main()
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/text", TextRequest).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}

type RequestBody struct {
    Text 	string   `json:"text,omitempty"`
}
type ResponseBody struct {
    Vocal  	string `json:"vox,omitempty"`
}


func TextRequest(w http.ResponseWriter, r *http.Request) {

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var req RequestBody
	err = json.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var nlu = c.UnderstandText(req.Text)
	var nlg = CallSkill(nlu)

	var resp ResponseBody
	resp.Vocal = nlg.Sentence
	json.NewEncoder(w).Encode(resp)
}

func CallSkill(nlu cerebro.NLU) anima.NLG {
	return anima.NLG{Sentence:"It is {time} o'clock"}
}
