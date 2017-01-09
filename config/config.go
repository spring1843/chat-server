package config

type Config struct {
	TelnetPort    int    `json:"telnet_port"`
	RestPort      int    `json:"rest_port"`
	WebsocketPort int    `json:"websocket_port"`
	LogFile       string `json:"log_file"`
	IP            string `json:"ip"`
}
