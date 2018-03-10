package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	node struct {
		auth string
		urls []string
		name string
	}
	web struct {
		listenAddress string `yaml:listen-address`
		metricURI     string `yaml:telemetry-path`
	}
}

func configure(filePath string) (config, error) {
	var c config
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
