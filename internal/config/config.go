package config

import (
	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
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

var config specification

func init() {
	err := envconfig.Process("borgmatic_exporter", &config)
	if err != nil {
		logs.Logger.Error(err.Error())
	}
}

func Load() (*specification, error) {
	return &config, nil
}
