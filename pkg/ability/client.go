package ability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.milobella.com/milobella/ability-sdk-go/pkg/ability"
	"gitlab.milobella.com/milobella/oratio/pkg/anima"
	"gitlab.milobella.com/milobella/oratio/pkg/cerebro"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	host string
	port int
	url string
	client http.Client
}

func NewClient(host string, port int) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{host: host, port: port, url: url, client: http.Client{}}
}

func (c Client) makeRequest(request ability.Request) (response ability.Response, err error) {
	endpoint := strings.Join([]string{c.url, "resolve", request.Nlu.BestIntent}, "/")
	postBody, err := json.Marshal(request)
	if err != nil {
		log.Print(err)
		return
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postBody))
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Print(err)
		return
	}
	return
}

func (c Client) CallAbility(nlu cerebro.NLU) (nlg anima.NLG, visu interface{}) {
	request := ability.Request{Nlu: nlu}
	result, err := c.makeRequest(request)
	if err != nil {
		log.Print(err)
		nlg.Sentence = "error"
		return
	}

	nlg = result.Nlg
	visu = result.Visu
	return
}