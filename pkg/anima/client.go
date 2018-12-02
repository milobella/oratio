package anima

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (c Client) makeRequest(nlg NLG) (result string, err error) {
	restituteEndpoint := strings.Join([]string{c.url, "restitute"}, "/")
	jsonNLG, err := json.Marshal(nlg)
	if err != nil {
		log.Print(err)
		return
	}
	req, err := http.NewRequest("POST", restituteEndpoint, bytes.NewBuffer(jsonNLG))
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

	return string(body), err
}

func (c Client) GenerateSentence(nlg NLG) (result string) {
	result, err := c.makeRequest(nlg)
	if err != nil {
		log.Print(err)
		result = "erreur"
	}
	return
}