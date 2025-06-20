package main

import (
	"encoding/base64"
	"fmt"
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

var isTesting bool

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: no .env file found.")
	}
	TV_IP := os.Getenv("TV_IP")
	if TV_IP == "" {
		log.Fatal("TV_IP environment variable is not set.")
	}
	CLIENT_PW := os.Getenv("CLIENT_PW")
	if CLIENT_PW == "" {
		fmt.Println("Warning: CLIENT_PW environment variable is not set. Gopher Remote will not reject unauthorized requests.")
	}
	APP_NAME := base64.StdEncoding.EncodeToString([]byte("Gopher Remote"))

	token := ""
	if data, err := os.ReadFile(TOKEN_FILE); err == nil {
		token = strings.TrimSpace(string(data))
	}

	socket := socket{
		ip:      TV_IP,
		port:    TV_PORT,
		appName: APP_NAME,
		token:   token,
	}
	defer socket.close()

	api := api{
		socket: &socket,
		pw:     CLIENT_PW,
	}

	http.HandleFunc("/increase-brightness", api.increaseBrightness)
	http.HandleFunc("/decrease-brightness", api.decreaseBrightness)

	fmt.Printf("Running  Gopher Remote on port %s\n", APP_PORT[1:])
	err = http.ListenAndServe(APP_PORT, nil)

	if err != nil {
		log.Fatal(err)
	}
}
