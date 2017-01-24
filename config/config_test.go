package config_test

import (
	"testing"

	"github.com/spring1843/chat-server/config"
)

func TestCanReadTestConfig(t *testing.T) {
	config := config.FromFile("./config_test.json")

	expectedIP := "6.6.6.6"
	if config.IP != expectedIP {
		t.Fatalf("Could not read from file, expected config value IP to be %s, got %s instead", expectedIP, config.IP)
	}
}
