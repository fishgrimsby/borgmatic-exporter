package config

import (
	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	Host      string `default:"0.0.0.0"`
	Port      string `default:"8090"`
	Endpoint  string `default:"metrics"`
	Config    string `default:""`
	Debug     bool   `default:"false"`
	LogFormat string `default:"keyvalue"`
}

func Load() (*specification, error) {
	var config specification
	err := envconfig.Process("borgmatic_exporter", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
