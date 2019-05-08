package anima

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
	host   string
	port   int
	url    string
	name   string
	client http.Client
}

func NewClient(host string, port int) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{host: host, port: port, url: url, client: http.Client{}, name: "anima"}
}

func (c Client) makeRequest(nlg NLG) (result string, err error) {
	restituteEndpoint := strings.Join([]string{c.url, "restitute"}, "/")
	jsonNLG, err := json.Marshal(nlg)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		return
	}
	req, err := http.NewRequest("POST", restituteEndpoint, bytes.NewBuffer(jsonNLG))
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

	return string(body), err
}

func (c Client) GenerateSentence(nlg NLG) (result string) {
	result, err := c.makeRequest(nlg)
	if err != nil {
		logrus.WithField("client", c.name).Error(err)
		result = "erreur"
	}
	return
}
