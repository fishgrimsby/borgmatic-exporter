package borg

import "os/exec"

// Allows unit testing exec.CommandContext output
var execCommand = exec.CommandContext
