package cerebro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	host   string
	port   int
	url    string
	client http.Client
}

func NewClient(host string, port int) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{host: host, port: port, url: url, client: http.Client{}}
}

func (c Client) UnderstandText(t string) (result NLU) {
	result, err := c.makeRequest(t)
	if err != nil {
		result.BestIntent = "error"
		return
	}

	c.bestNLU(&result)
	return
}

func (c Client) bestNLU(result *NLU) {
	var bestScore float32 = 0
	for _, intent := range result.Intents {
		if intent.Score > bestScore {
			result.BestIntent = intent.Label
			bestScore = intent.Score
		}
	}
}

func (c Client) makeRequest(query string) (result NLU, err error) {
	understandEndpoint := strings.Join([]string{c.url, "understand"}, "/")
	req, err := http.NewRequest("GET", understandEndpoint, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Error(err)
		return
	}

	err = json.Unmarshal(body, &result)
	return
}
