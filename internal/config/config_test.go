package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	want := &specification{}
	got, err := Load()

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if reflect.TypeOf(want) != reflect.TypeOf(got) {
		t.Fatalf(`Load() = %v, want match for %#q, nil`, reflect.TypeOf(got), reflect.TypeOf(want))
	}

}
