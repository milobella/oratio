package cerebro

import (
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

func (c Client) UnderstandText(t string) (nlu NLU) {
	result, err := c.makeRequest(t)
	if err != nil {
		nlu.Intent = "error"
		return
	}

	nlu.Intent = c.interpretIntent(result)
	return
}

func (c Client) interpretIntent(result map[string]float32) (best string) {
	var bestScore float32 = 0
	for category, score := range result {
		if score > bestScore {
			best = category
			bestScore = score
		}
	}
	return
}

func (c Client) makeRequest(query string) (result map[string]float32, err error) {
	understandEndpoint := strings.Join([]string{c.url, "understand"}, "/")
	req, err := http.NewRequest("GET", understandEndpoint, nil)
	if err != nil {
		log.Print(err)
		return
	}
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

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

	err = json.Unmarshal(body, &result)
	return
}
