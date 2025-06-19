package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type tvResponse struct {
	Event string         `json:"event"`
	Data  tvResponseData `json:"data"`
}
type tvResponseData struct {
	Token string `json:"token"`
}

type socket struct {
	ip         string
	port       string
	appName    string
	token      string
	connection *websocket.Conn
	testing    bool
	tokenFile  string
}

const DEFAULT_TOKEN_FILE = ".tv_token"

func (s *socket) connect() error {
	if s.connection != nil {
		return nil
	}

	protocol := "wss"
	if s.testing {
		protocol = "ws"
	}
	wsUrl := fmt.Sprintf(
		"%s://%s:%s/api/v2/channels/samsung.remote.control?name=%s",
		protocol, s.ip, s.port, s.appName,
	)
	if s.token != "" {
		wsUrl += "&token=" + s.token
	}
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		return fmt.Errorf("could not connect to address: %w", err)
	}

	connChn := make(chan error, 1)

	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				connChn <- err
				return
			}

			var msg tvResponse
			if err := json.Unmarshal(data, &msg); err != nil {
				connChn <- err
				return
			}

			if msg.Event == "ms.channel.connect" {
				fmt.Println("Connected to TV")
				respToken := msg.Data.Token
				if s.token != respToken {
					s.token = respToken
					if s.tokenFile == "" {
						s.tokenFile = DEFAULT_TOKEN_FILE
					}
					err = os.WriteFile(s.tokenFile, []byte(respToken), 0644)
					if err != nil {
						connChn <- fmt.Errorf("error: could not write token to file: %w", err)
						continue
					}
					connChn <- nil
					return
				}
				connChn <- nil
				return
			}
		}
	}()

	err = <-connChn
	if err != nil {
		conn.Close()
		return fmt.Errorf("could not connect to tv: %w", err)
	}

	s.connection = conn
	return nil
}

func (s *socket) close() {
	if s.connection != nil {
		err := s.connection.Close()
		if err != nil {
			log.Println(err)
		}
		s.connection = nil
	}
}
