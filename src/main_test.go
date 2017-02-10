package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/bootstrap"
)

const configFile = "config.json"

// TestDefaultConfig runs the application using default config to ensure it runs
func TestDefaultConfig(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		config := config.FromFile(configFile)
		bootstrap.NewBootstrap(config)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestDefaultConfig")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	if err != nil {
		t.Fatalf("Couldnt run application with values in config.json, Test exited with status 1, expected 0. Error %s", err)
	}
}

// TestBadConfig runs the application using a bad config config to ensure it ends with error status 1
func TestBadConfig(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		config := config.FromFile(configFile)

		// Make config invalid
		config.TelnetAddress = "-1"

		bootstrap.NewBootstrap(config)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestBadConfig")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Did not fail run with invalid config.json values, Test exited with status 0, expected 1")
}
