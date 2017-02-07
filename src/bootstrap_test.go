package main

import (
	"os"
	"os/exec"
	"testing"

	"net/http"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

func TestCanStartWebWithHTTP(t *testing.T) {
	config := config.FromFile("./config.json")
	config.WebAddress += "1"
	srv := getTLSServer(getmux(), config.WebAddress)
	go func() {
		t.Logf("Starting http on %s", config.WebAddress)
		if err := srv.ListenAndServe(); err != nil {
			logs.FatalIfErrf(err, "Couldn't start http on %s", config.WebAddress)
		}
	}()

	resp, err := http.Get("http://" + config.WebAddress + "/api/status")
	if err != nil {
		t.Fatalf("couldn't get api status endpoint after serving http. Error %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestCanStartWebWithHTTPS(t *testing.T) {
	config := config.FromFile("./config.json")
	config.WebAddress += "2"
	srv := getTLSServer(getmux(), config.WebAddress)
	go func() {
		t.Logf("Starting https on %s", config.WebAddress)
		if err := srv.ListenAndServeTLS("tls.crt", "tls.key"); err != nil {
			logs.FatalIfErrf(err, "Couldn't start https on %s", config.WebAddress)
		}
	}()

	resp, err := http.Get("https://" + config.WebAddress + "/api/status")
	if err != nil {
		t.Fatalf("couldn't get api status endpoint after serving https. Error %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestEmptyDrivers(t *testing.T) {
	config := new(config.Config)
	bootstrap(*config)
	if chatServer == nil {
		t.Fatalf("Empty web and telnet addresses did not start the server.")
	}
}

func TestErrorOnInvalidTelnet(t *testing.T) {
	config := new(config.Config)
	config.TelnetAddress = "-1"
	if err := startTelnet(*config); err == nil {
		t.Fatalf("Expected error on invalid telnet address")
	}
}

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
		config.TelnetAddress = "-1"

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
