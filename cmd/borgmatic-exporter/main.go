package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fishgrimsby/borgmatic-exporter/internal/config"
	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
	"github.com/fishgrimsby/borgmatic-exporter/internal/metrics"
)

func main() {
	config, err := config.Load()
	logs.Configure(config.Debug, config.LogFormat)

	if err != nil {
		logs.Logger.Error(err.Error())
		os.Exit(1)
	}

	collector := metrics.New(config.Config)
	prometheus.MustRegister(collector)

	http.Handle(fmt.Sprintf("/%s", config.Endpoint), promhttp.Handler())
	addr := config.Host + ":" + config.Port
	logs.Logger.Info(fmt.Sprintf("listening on http://%s/%s", addr, config.Endpoint))
	err = http.ListenAndServe(addr, nil)

	if errors.Is(err, http.ErrServerClosed) {
		logs.Logger.Error("server closed",
			"error", err.Error())
	} else if err != nil {
		logs.Logger.Error("error starting server",
			"error", err.Error())
		logs.Logger.Error("exiting application")
		os.Exit(1)
	}
}
