package bootstrap

import (
	"testing"

	"net/http"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const configFile = "../../config.json"

func TestCanStartWebWithHTTP(t *testing.T) {
	t.Skipf("Doesnt start on build server.")
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
