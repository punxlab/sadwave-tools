package client

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type SadwaveClient struct {
	host string
}

func NewSadwaveClient(host string) *SadwaveClient {
	return &SadwaveClient{host: host}
}

func (c *SadwaveClient) SetCustomPhoto(eventUrl *url.URL, photoUrl *url.URL) error {
	body, err := json.Marshal(map[string]string{
		"event": eventUrl.String(),
		"photo": photoUrl.String(),
	})

	if err != nil {
		return errors.Wrap(err, "marshal request body")
	}

	reqUrl := c.host + "/api/photos"
	resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "post request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected status %s", resp.Status)
	}

	return nil
}
