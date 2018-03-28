package main

import (
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
	"os"
)

type config struct {
	Node struct {
		Auth struct {
			User     string `yaml:"-"`
			Password string `yaml:"-"`
		}
		URLs    []string
		Name    string
		Refresh time.Duration
	}
	Web struct {
		Adress string
		URI    string
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
	if len(os.Getenv("CBE_USER")) > 0 {
		c.Node.Auth.User = os.Getenv("CBE_USER")
	}
	if len(os.Getenv("CBE_PASSWORD")) > 0 {
		c.Node.Auth.Password = os.Getenv("CBE_PASSWORD")
	}
	if len(c.Node.URLs) == 0 {
		return c, fmt.Errorf("Can't find cluster urls in config\n")
	}
	//defaults
	if c.Web.Adress == "" {
		c.Web.Adress = ":9131"
	}
	if c.Web.URI == "" {
		c.Web.URI = "/metrics"
	}
	if c.Node.Refresh == 0 {
		c.Node.Refresh = 5 * time.Second
	}
	return c, nil
}
