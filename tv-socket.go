package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func getSocket(ip, port, appName, token string) (*websocket.Conn, error) {
	wsUrl := fmt.Sprintf(
		"wss://%s:%s/api/v2/channels/samsung.remote.control?name=%s",
		ip, port, appName,
	)
	if token != "" {
		wsUrl += "&token=" + token
	}

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			var msg struct {
				Data struct {
					Token string `json:"token"`
				} `json:"data"`
				Event string `json:"event"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Fatal(err)
			}

			if msg.Event == "ms.channel.connect" {
				fmt.Println("Connected to TV")
				if token != "" {
					continue
				}
				err = os.WriteFile(TOKEN_FILE, []byte(msg.Data.Token), 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return conn, nil
}
