package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("BORGMATIC_EXPORTER_HOST", "127.0.0.1")
	t.Setenv("BORGMATIC_EXPORTER_PORT", "9090")
	t.Setenv("BORGMATIC_EXPORTER_ENDPOINT", "prom")

	got, err := Load()

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if got.Host != "127.0.0.1" {
		t.Errorf("Expected Host 127.0.0.1, got %q", got.Host)
	}
	if got.Port != "9090" {
		t.Errorf("Expected Port 9090, got %q", got.Port)
	}
	if got.Endpoint != "prom" {
		t.Errorf("Expected Endpoint prom, got %q", got.Endpoint)
	}
}
