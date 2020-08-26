package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
)

const configPath = "config.json"

type Config struct {
	API           string `json:"api"`
	AdminUser     string `json:"user"`
	AdminPassword string `json:"password"`
	SourceCSVPath string `json:"source"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	if cfg.API == "" {
		return nil, errors.New("empty api url")
	}

	_, err = url.Parse(cfg.API)
	if err != nil {
		return nil, errors.Wrap(err, "invalid api url")
	}

	if cfg.SourceCSVPath == "" {
		return nil, errors.New("empty source file path")
	}

	if cfg.AdminUser == "" {
		return nil, errors.New("empty user")
	}

	if cfg.AdminPassword == "" {
		return nil, errors.New("empty password")
	}

	return cfg, nil
}
