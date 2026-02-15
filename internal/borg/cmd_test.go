package borg

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestExecCommand(t *testing.T) {
	want := exec.CommandContext
	got := execCommand

	if reflect.TypeOf(want) != reflect.TypeOf(got) {
		t.Fatalf(`execCommand = %v, want match for %#q, nil`, reflect.TypeOf(got), reflect.TypeOf(want))
	}
}
