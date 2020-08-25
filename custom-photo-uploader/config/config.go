package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
)

const configPath = "./config.json"

type Config struct {
	API           *url.URL
	AdminUser     string
	AdminPassword string
	SourceCSVPath string
}

func NewConfig() (*Config, error) {
	var raw = &struct {
		API           string `json:"api"`
		AdminUser     string `json:"user"`
		AdminPassword string `json:"password"`
		SourceCSVPath string `json:"api"`
	}{}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	err = json.Unmarshal(data, raw)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	if raw.API == "" {
		return nil, errors.Wrap(err, "empty api url")
	}

	api, err := url.Parse(raw.API)
	if err != nil {
		return nil, errors.Wrap(err, "invalid api url")
	}

	if raw.SourceCSVPath == "" {
		return nil, errors.New("empty source file path")
	}

	if raw.AdminUser == "" {
		return nil, errors.New("empty user")
	}

	if raw.AdminPassword == "" {
		return nil, errors.New("empty password")
	}

	return &Config{
		API:           api,
		AdminUser:     raw.AdminUser,
		AdminPassword: raw.AdminPassword,
		SourceCSVPath: raw.SourceCSVPath,
	}, nil
}
