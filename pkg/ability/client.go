package ability

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.milobella.com/milobella/ability-sdk-go/pkg/ability"
	"gitlab.milobella.com/milobella/oratio/pkg/anima"
	"gitlab.milobella.com/milobella/oratio/pkg/cerebro"
)

// Client : Ability HTTP client
type Client struct {
	host   string
	port   int
	url    string
	client http.Client
	name   string
}

// NewClient : ctor
func NewClient(host string, port int, name string) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{host: host, port: port, url: url, client: http.Client{}, name: name}
}

func (c Client) makeRequest(request ability.Request) (response ability.Response, err error) {
	endpoint := strings.Join([]string{c.url, "resolve", request.Nlu.BestIntent}, "/")
	postBody, err := json.Marshal(request)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postBody))
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

	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}
	return
}

// CallAbility : Requests the ability
func (c Client) CallAbility(nlu cerebro.NLU) (nlg anima.NLG, visu interface{}, autoReprompt bool) {
	// By default the auto reprompt is false
	autoReprompt = false
	request := ability.Request{Nlu: nlu}
	result, err := c.makeRequest(request)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		nlg.Sentence = "error"
		return
	}

	nlg = result.Nlg
	visu = result.Visu
	autoReprompt = result.AutoReprompt
	return
}
