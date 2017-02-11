package bootstrap

import (
	"net/http"
	"os"
	"testing"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const configFile = "../../config.json"

func TestCanStartWebWithHTTP(t *testing.T) {
	if os.Getenv("SKIP_NETWORK") == "1" {
		t.Skipf("Skipping test SKIP_NETWORK set to %q", os.Getenv("SKIPNETWORK"))
	}
	config := config.FromFile(configFile)
	config.WebAddress += "1"
	srv := getTLSServer(getMultiplexer(), config.WebAddress)
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

func TestEmptyDrivers(t *testing.T) {
	config := new(config.Config)
	NewBootstrap(*config)
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
