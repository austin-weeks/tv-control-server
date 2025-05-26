package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	TV_PORT    = "8002"
	TOKEN_FILE = ".tv_token"
	APP_PORT   = ":1234"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	TV_IP := os.Getenv("TV_IP")
	CLIENT_PW := os.Getenv("CLIENT_PW")
	APP_NAME := base64.StdEncoding.EncodeToString([]byte("Gopher Remote"))

	token := ""
	// Read token if file exists
	if _, err := os.Stat(TOKEN_FILE); err == nil || !os.IsNotExist(err) {
		if data, err := os.ReadFile(TOKEN_FILE); err == nil {
			token = strings.TrimSpace(string(data))
		}
	}

	socket, err := getSocket(TV_IP, TV_PORT, APP_NAME, token)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	api := api{
		ws: socket,
		pw: CLIENT_PW,
	}

	http.HandleFunc("/increase-brightness", api.increaseBrightness)
	http.HandleFunc("/decrease-brightness", api.decreaseBrightness)
	err = http.ListenAndServe(APP_PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
