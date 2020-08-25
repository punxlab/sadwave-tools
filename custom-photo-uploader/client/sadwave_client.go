package client

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"path"
)

type SadwaveClient struct {
	url *url.URL
}

func NewSadwaveClient(url *url.URL) *SadwaveClient {
	return &SadwaveClient{url: url}
}

func (c *SadwaveClient) SetCustomPhoto(eventUrl *url.URL, photoUrl *url.URL) error {
	body, err := json.Marshal(map[string]string{
		"event": eventUrl.String(),
		"photo": photoUrl.String(),
	})

	if err != nil {
		return errors.Wrap(err, "marshal request body")
	}

	resp, err := http.Post(path.Join(c.url.Path, "api/photos"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "post request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected status %s", resp.Status)
	}

	return nil
}
