package ability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/milobella/oratio/pkg/anima"
	"github.com/sirupsen/logrus"
)

// Client : Ability HTTP client
type Client struct {
	Host   string
	Port   int
	url    string
	client http.Client
	Name   string
}

// NewClient : ctor
func NewClient(host string, port int, name string) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{Host: host, Port: port, url: url, client: http.Client{}, Name: name}
}

func (c Client) makeRequest(request Request) (response Response, err error) {
	endpoint := strings.Join([]string{c.url, "resolve"}, "/")
	postBody, err := json.Marshal(request)
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		return
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postBody))
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		return
	}

	logrus.WithField("client", c.Name).WithField("status", resp.StatusCode).Infof("%s %s", req.Method, req.URL)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		return
	}
	return
}

// CallAbility : Requests the ability
func (c Client) CallAbility(request Request) (nlg anima.NLG, visu interface{}, autoReprompt bool, context Context) {
	// By default the auto reprompt is false
	autoReprompt = false
	result, err := c.makeRequest(request)
	if err != nil {
		logrus.WithField("client", c.Name).Error(err)
		nlg.Sentence = "error"
		return
	}

	nlg = result.Nlg
	visu = result.Visu
	autoReprompt = result.AutoReprompt
	context = result.Context
	return
}
