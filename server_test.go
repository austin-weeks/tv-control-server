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
	mux.HandleFunc(
		"/api/v2/channels/samsung.remote.control",
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			msg := tvResponse{
				Event: "ms.channel.connect",
				Data:  tvResponseData{Token: token},
			}
			b, err := json.Marshal(msg)
			if err != nil {
				t.Fatal(err)
			}
			conn.WriteMessage(websocket.TextMessage, b)

			// Check Sent Keys
			keyInd := 0
			for {
				_, msgBytes, err := conn.ReadMessage()
				if err != nil {
					t.Error("Error reading client message:", err)
				}
				if expectedMacros == nil {
					continue
				}
				var msg keyMsg
				err = json.Unmarshal(msgBytes, &msg)
				if err != nil {
					t.Error("Could not deserialize client message:", err)
				}
				sentKey := msg.Params.DataOfCmd
				expectedKey := expectedMacros[keyInd].key
				if sentKey != expectedKey {
					t.Fatalf("Sent key %s does not match expected %s", sentKey, expectedKey)
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
			t.Error("Test server failed to serve", err)
		}
	}()

	return func() {
		server.Close()
		ln.Close()
	}
}
