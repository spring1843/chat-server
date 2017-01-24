package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/spring1843/chat-server/config"
)

// TestCanRunDefaultConfig run bootstrap and expects it to be able to run with config.json values
func TestCanRunDefaultConfig(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		config := config.FromFile("./config.json")
		bootstrap(config)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCanRunDefaultConfig")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	if err != nil {
		t.Fatalf("Couldnt run with config.json values, Test exited with status 1, expected 1. Error %s", err)
	}
}

// TestCanCrashOnBadConfig run bootstrap and expects it to exit with status 1 (fatal error) when config values are not valid
func TestCanCrashOnBadConfig(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		config := config.FromFile("./config.json")

		// Make config invalid
		config.TelnetPort = -1

		bootstrap(config)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCanCrashOnBadConfig")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Did not fail run with invalid config.json values, Test exited with status 0, expected 1")
}
