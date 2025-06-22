package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var isTesting bool

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	APP_NAME := base64.StdEncoding.EncodeToString([]byte(config.AppName))

	token := ""
	if data, err := os.ReadFile(config.TokenFile); err == nil {
		token = strings.TrimSpace(string(data))
	}

	socket := socket{
		ip:      config.TvIP,
		port:    config.TvPort,
		appName: APP_NAME,
		token:   token,
	}
	defer socket.close()

	api := api{
		socket:             &socket,
		pw:                 config.ClientPassword,
		brightnessPosition: config.BrightnessLocation,
		initialDelay:       time.Duration(config.InitialDelay) * time.Millisecond,
	}

	http.HandleFunc("/increase-brightness", api.increaseBrightness)
	http.HandleFunc("/decrease-brightness", api.decreaseBrightness)

	fmt.Printf(" Running %s on port %d\n", config.AppName, config.AppPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.AppPort), nil)
	log.Fatal(err)
}
