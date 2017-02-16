package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config holds the values ChatServer needs in order to run
type Config struct {
	TelnetAddress string `json:"telnet_address"`
	WebAddress    string `json:"web_address"`
	StaticWeb     string `json:"static_web"`
	HTTPS         bool   `json:"https"`
}

// FromFile parses Config from a .json file
func FromFile(configFile string) Config {
	fileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config := new(Config)
	err = json.Unmarshal([]byte(fileContents), &config)
	if err != nil {
		log.Fatalf("Error parsing JSON config file: %s", err)
	}

	return *config
}
