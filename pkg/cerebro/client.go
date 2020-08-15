package cerebro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	host               string
	port               int
	url                string
	client             http.Client
	name               string
	understandEndpoint string
}

// buildEndpoint: Ensure that endpoint start with /
func buildEndpoint(endpoint string) string {
	if !strings.HasPrefix(endpoint, "/") {
		return "/" + endpoint
	}
	return endpoint
}

func NewClient(host string, port int, understandEndpoint string) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{
		host:   host,
		port:   port,
		url:    url,
		client: http.Client{},
		name:   "cerebro",
		understandEndpoint: buildEndpoint(understandEndpoint),
	}
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
	understandEndpoint := c.url + c.understandEndpoint
	reqBody := []byte(fmt.Sprintf("{\"text\": \"%s\"}", query))
	req, err := http.NewRequest("POST", understandEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}

	logrus.WithField("client", c.name).WithField("status", resp.StatusCode).Infof("%s %s", req.Method, req.URL)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}

	logrus.WithField("client", c.name).Debug(string(body))

	err = json.Unmarshal(body, &result)
	return
}
