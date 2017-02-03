package config_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/config"
)

func TestCanReadTestConfig(t *testing.T) {
	config := config.FromFile("./config_test.json")

	expectedAddress := "6.6.6.6:4004"
	if config.WebAddress != expectedAddress {
		t.Fatalf("Could not read from file, expected config value IP to be %s, got %s instead", expectedAddress, config.WebAddress)
	}
}
