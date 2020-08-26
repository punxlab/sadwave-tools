package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type SadwaveClient struct {
	host string
	cl   *http.Client
}

func NewSadwaveClient(host string) *SadwaveClient {
	return &SadwaveClient{host: host, cl: new(http.Client)}
}

func (c *SadwaveClient) SetCustomPhoto(eventUrl *url.URL, photoUrl *url.URL, authToken string) error {
	body, err := json.Marshal(map[string]string{
		"event": eventUrl.String(),
		"photo": photoUrl.String(),
	})

	if err != nil {
		return errors.Wrap(err, "marshal request body")
	}

	req, err := http.NewRequest("POST", c.host+"/api/photos", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "new request")
	}

	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cl.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected status %s", resp.Status)
	}

	return nil
}

func (c *SadwaveClient) RequestToken(login, password string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})

	if err != nil {
		return "", errors.Wrap(err, "marshal request body")
	}

	req, err := http.NewRequest("POST", c.host+"/api/authentication", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cl.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.Errorf("unexpected status %s", resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	res := struct {
		Token string `json:"token"`
	}{}

	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal response body")
	}

	return res.Token, nil
}
