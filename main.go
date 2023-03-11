package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fishgrimsby/borgmatic-exporter/internal/metrics"
)

type Specification struct {
	Port     string `default:"8090"`
	Endpoint string `default:"metrics"`
	Config   string `default:""`
}

func main() {
	var s Specification
	err := envconfig.Process("borgmatic_exporter", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	collector := metrics.New(s.Config)
	prometheus.MustRegister(collector)

	http.Handle(fmt.Sprintf("/%s", s.Endpoint), promhttp.Handler())

	fmt.Printf("Listening on port %s with endpoint /%s", s.Port, s.Endpoint)
	err = http.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
