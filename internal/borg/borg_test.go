package borg

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	want := &borg{}
	got, err := New()

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if reflect.TypeOf(want) != reflect.TypeOf(got) {
		t.Fatalf(`New() = %v, want match for %#q, nil`, reflect.TypeOf(got), reflect.TypeOf(want))
	}
}
