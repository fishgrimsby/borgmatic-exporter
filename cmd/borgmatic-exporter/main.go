package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fishgrimsby/borgmatic-exporter/internal/config"
	"github.com/fishgrimsby/borgmatic-exporter/internal/logs"
	"github.com/fishgrimsby/borgmatic-exporter/internal/metrics"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %s\n", err)
		os.Exit(1)
	}

	logs.Configure(config.Debug, config.LogFormat)

	collector := metrics.New(config.Config, config.Timeout)
	prometheus.MustRegister(collector)

	http.Handle(fmt.Sprintf("/%s", config.Endpoint), promhttp.Handler())
	addr := config.Host + ":" + config.Port

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{Addr: addr}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logs.Logger.Error("server shutdown error", "error", err.Error())
		}
	}()

	logs.Logger.Info(fmt.Sprintf("listening on http://%s/%s", addr, config.Endpoint))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logs.Logger.Error("error starting server", "error", err.Error())
		os.Exit(1)
	}
}
