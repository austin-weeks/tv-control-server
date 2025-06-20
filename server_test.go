package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

const TEST_IP = "127.0.0.1"
const TEST_PORT = "12345"
const TEST_TOKEN = "test-token"

func getTestSocket() (socket, error) {
	socket := socket{
		ip:    TEST_IP,
		port:  TEST_PORT,
		token: TEST_TOKEN,
	}
	err := socket.connect()
	return socket, err
}

func startTestWSServer(t *testing.T, token string, expectedMacros []macro) func() {
	upgrader := websocket.Upgrader{}
	server := http.Server{}

	mux := http.NewServeMux()
	server.Handler = mux

	errCh := make(chan error, 1)
	mux.HandleFunc(
		"/api/v2/channels/samsung.remote.control",
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				errCh <- err
				return
			}
			defer conn.Close()
			msg := tvResponse{
				Event: "ms.channel.connect",
				Data:  tvResponseData{Token: token},
			}
			b, err := json.Marshal(msg)
			if err != nil {
				errCh <- err
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				errCh <- err
				return
			}

			// Check Sent Keys
			keyInd := 0
			for {
				_, msgBytes, err := conn.ReadMessage()
				if err != nil {
					errCh <- fmt.Errorf("error reading client message: %w", err)
					return
				}
				if expectedMacros == nil {
					continue
				}
				var msg keyMsg
				err = json.Unmarshal(msgBytes, &msg)
				if err != nil {
					errCh <- fmt.Errorf("could not deserialize client message: %w", err)
					return
				}
				sentKey := msg.Params.DataOfCmd
				expectedKey := expectedMacros[keyInd].key
				if sentKey != expectedKey {
					errCh <- fmt.Errorf("sent key %s does not match expected %s", sentKey, expectedKey)
					return
				}
				keyInd++
			}
		},
	)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", TEST_IP, TEST_PORT))
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		err := server.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("test server failed to serve: %w", err)
		}
	}()

	return func() {
		server.Close()
		ln.Close()
		select {
		case err := <-errCh:
			if err != nil {
				t.Error("Test server error:", err)
			}
		default:
		}
	}
}
