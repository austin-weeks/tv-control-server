package main

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	configPath = t.TempDir() + "test_config.json"
	clearConfigFile := func() {
		if _, err := os.Stat(configPath); err == nil {
			err := os.Remove(configPath)
			if err != nil {
				panic(err)
			}
		}
	}

	t.Run("Reads and parses config file", func(t *testing.T) {
		clearConfigFile()
		err := os.WriteFile(configPath, []byte(`{
			"app_name": "test-name",
			"app_port": "test-app-port",
			"token_file": "test-token-file",
			"tv_port": "test-tv-port",
			"tv_ip": "test-tv-ip",
			"client_password": "test-client-password",
			"brightness_location": 1,
			"initial_delay_ms": 2000
		}`), 0644)
		if err != nil {
			panic(err)
		}
		config, err := getConfig()
		if err != nil {
			t.Fatal("Failed to get config:", err)
		}
		if config.AppName != "test-name" ||
			config.AppPort != "test-app-port" ||
			config.TokenFile != "test-token-file" ||
			config.TvPort != "test-tv-port" ||
			config.TvIP != "test-tv-ip" ||
			config.ClientPassword != "test-client-password" ||
			config.BrightnessLocation != 1 ||
			config.InitialDelay != 2000 {
			t.Errorf("Config, %v, does not contain expected values", config)
		}
	})

	t.Run("Fails if no config file", func(t *testing.T) {
		clearConfigFile()
		_, err := getConfig()
		if err == nil {
			t.Error("Expected failure with missing config file")
		}
	})

	t.Run("Fails if no TV IP address", func(t *testing.T) {
		clearConfigFile()
		err := os.WriteFile(configPath, []byte("{}"), 0644)
		if err != nil {
			panic(err)
		}
		_, err = getConfig()
		if err == nil {
			t.Error("Expected failure with missing TV IP address")
		}
	})

	t.Run("Sets default values", func(t *testing.T) {
		clearConfigFile()
		err := os.WriteFile(configPath, []byte(`{
			"tv_ip": "test-tv-ip"
		}`), 0644)
		if err != nil {
			panic(err)
		}
		config, err := getConfig()
		if err != nil {
			t.Fatal("Failed to get config:", err)
		}
		if config.TvIP != "test-tv-ip" ||
			config.AppName != DEFAULT_APP_NAME ||
			config.AppPort != DEFAULT_APP_PORT ||
			config.TvPort != DEFAULT_TV_PORT ||
			config.TokenFile != DEFAULT_TOKEN_FILE ||
			config.ClientPassword != "" ||
			config.BrightnessLocation != DEFAULT_BRIGHT_LOC ||
			config.InitialDelay != DEFAULT_INITIAL_DELAY {
			t.Errorf("Config, %v, does not contain default values", config)
		}
	})

	clearConfigFile()
}
