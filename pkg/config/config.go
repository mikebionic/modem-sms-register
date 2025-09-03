package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SerialPort string `json:"serial_port"`
	Hostname   string `json:"hostname"`
	URLPrefix  string `json:"url_prefix"`
	URLPath    string `json:"url_path"`
	SHAKey     string `json:"sha_key"`
	Baud       int    `json:"baud,omitempty"`
	Timeout    string `json:"timeout,omitempty"`
	Period     string `json:"period,omitempty"`
}

func ReadConfig(source string) (*Config, error) {
	raw, err := os.ReadFile(source)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config file")
	}

	var c Config
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, errors.Wrap(err, "error parsing config JSON")
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) validate() error {
	if c.SerialPort == "" {
		return errors.New("serial_port is required")
	}
	if c.Hostname == "" {
		return errors.New("hostname is required")
	}
	if c.SHAKey == "" {
		return errors.New("sha_key is required")
	}
	return nil
}

func GetConfigData() (*Config, error) {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config.json"
	}
	cfg, err := ReadConfig(configFile)
	if err != nil {
		log.WithError(err).WithField("config-file", configFile).Error("error loading configuration")
		return nil, err
	}
	return cfg, nil
}

func (c *Config) GetURL() string {
	return c.Hostname + c.URLPrefix + c.URLPath
}
