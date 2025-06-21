package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	AppName            string `json:"app_name"`
	AppPort            string `json:"app_port"`
	TokenFile          string `json:"token_file"`
	TvPort             string `json:"tv_port"`
	TvIP               string `json:"tv_ip"`
	ClientPassword     string `json:"client_password"`
	BrightnessLocation int    `json:"brightness_location"`
}

const (
	DEFAULT_APP_NAME   = "Gopher Remote"
	DEFAULT_APP_PORT   = "1234"
	DEFAULT_TV_PORT    = "8002"
	DEFAULT_TOKEN_FILE = ".tv_token"
	DEFAULT_BRIGHT_LOC = 3
)

var configPath = "config.json"

func getConfig() (config, error) {
	var config config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("could not read config.json: %w", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not deserialize config.json: %w", err)
	}
	if config.TvIP == "" {
		return config, fmt.Errorf("tv ip address was not provided in config.json")
	}
	if config.AppName == "" {
		config.AppName = DEFAULT_APP_NAME
	}
	if config.AppPort == "" {
		config.AppPort = DEFAULT_APP_PORT
	}
	if config.TokenFile == "" {
		config.TokenFile = DEFAULT_TOKEN_FILE
	}
	if config.TvPort == "" {
		config.TvPort = DEFAULT_TV_PORT
	}
	if config.BrightnessLocation == 0 {
		config.BrightnessLocation = DEFAULT_BRIGHT_LOC
	}
	return config, nil
}
