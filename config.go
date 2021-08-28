package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Serial_port string `json:"serial_port"`
	Hostname    string `json:"hostname"`
	URL_prefix  string `json:"url_prefix"`
	URL_Path    string `json:"url_path"`
	SHA_KEY     string `json:"sha_key"`
}

func ReadConfig(source string) (c *config, err error) {
	var raw []byte
	raw, err = ioutil.ReadFile(source)
	if err != nil {
		eMsg := "error reading config from file"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		eMsg := "error parsing config from json"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		c = nil
	}
	return
}

func get_config_data() (c *config, err error) {
	configFile := "config.json"
	c, err = ReadConfig(configFile)
	if err != nil {
		log.WithError(err).WithField("config-file", configFile).Error("error loading configuration")
		return
	}
	return
}

func get_port_and_url_from_config() (serial_port string, url_address string, err error) {
	conf, err := get_config_data()
	if err != nil {
		return
	}

	serial_port = conf.Serial_port
	url_address = fmt.Sprintf("%s%s%s", conf.Hostname, conf.URL_prefix, conf.URL_Path)

	return
}

func get_sha_key_from_config() (sha_key string, err error) {
	conf, err := get_config_data()
	if err != nil {
		return
	}
	sha_key = conf.SHA_KEY
	return
}
