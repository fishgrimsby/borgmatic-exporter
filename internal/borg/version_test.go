package borg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

// Test version is returned if Borg is installed
func TestGetVersionInstalled(t *testing.T) {
	execCommand = fakeExecCommandInstalled
	defer func() { execCommand = exec.CommandContext }()

	want := regexp.MustCompile(`(?m)^borg (?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	got, err := getVersion(context.Background())

	if err != nil {
		t.Fatalf("Expected nil error, got %#v", err)
	}

	if !want.MatchString(got) {
		t.Fatalf(`getVersion() = %v, want match for %#q, nil`, got, want)
	}
}

// Test for error if Borg is not installed
func TestGetVersionNotInstalled(t *testing.T) {
	execCommand = fakeExecCommandNotInstalled
	defer func() { execCommand = exec.CommandContext }()

	want := errors.New("executable not found")
	_, got := getVersion(context.Background())

	if got == nil {
		t.Fatalf(`getVersion() = %v, want match for, got %v`, want, got)
	}
}

// Test helpers
func TestGetVersionNotInstalledHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	os.Exit(1)
}

const borgVersionCmdResultInstalled = "borg 1.1.18"

func TestGetVersionInstalledHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, borgVersionCmdResultInstalled)
	os.Exit(0)
}

func fakeExecCommandInstalled(_ context.Context, command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGetVersionInstalledHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeExecCommandNotInstalled(_ context.Context, command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestGetVersionNotInstalledHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
