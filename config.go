package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

func configure(command string) (bool, error) {
	if command == "init" {
		flag.Parse()
	}

	if _, err := toml.DecodeFile(*ConfigFileName, &Configuration); err != nil {
		log.Fatal("configure", command, err)
		return false, err
	}

	return true, nil
}

// Config type include config sections
type Config struct {
	Core CoreConfig
}

// CoreConfig configuration server
type CoreConfig struct {
	IP              string
	Port            string
	NodeURL         []string `toml:"nodeurl"`
	RefreshInterval duration `toml:"refreshinterval"`
	Username        string
	Password        string
	MetricUri       string `toml:"metricUri""`
}

func (r *CoreConfig) getAddress() string {
	if r.IP == "" {
		r.IP = "0.0.0.0"
	}
	if r.Port == "" {
		r.Port = "9131"
	}
	return r.IP + ":" + r.Port
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
