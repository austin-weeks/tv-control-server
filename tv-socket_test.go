package main

import (
	"os"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {
	t.Run("Exercise Happy Path", func(t *testing.T) {
		token := "test-token"
		close := startTestWSServer(t, token, nil)
		defer close()
		socket := socket{
			ip:      TEST_IP,
			port:    TEST_PORT,
			testing: true,
			token:   token,
		}
		err := socket.connect()
		if err != nil {
			t.Error(err)
		}
	})

	for _, testToken := range []string{"", "invalid-token"} {
		t.Run("Token is written to file "+testToken, func(t *testing.T) {
			correctToken := "test-token"
			close := startTestWSServer(t, correctToken, nil)
			defer close()
			tmpFile := t.TempDir() + ".tv_token"
			socket := socket{
				ip:        TEST_IP,
				port:      TEST_PORT,
				testing:   true,
				token:     testToken,
				tokenFile: tmpFile,
			}
			err := socket.connect()
			if err != nil {
				t.Error(err)
			}
			data, err := os.ReadFile(tmpFile)
			if err != nil {
				t.Error("Token file not found:", err)
			}
			writtenToken := string(data)
			if writtenToken != correctToken {
				t.Errorf("Expected token %s does not match written token: %s", correctToken, writtenToken)
			}
		})
	}

	t.Run("Socket does not reconnect when connect", func(t *testing.T) {
		connection := websocket.Conn{}
		socket := socket{
			connection: &connection,
		}
		err := socket.connect()
		if err != nil {
			t.Errorf("Error while connecting socket: %s", err.Error())
		}
		if socket.connection != &connection {
			t.Error("Connection is not preserver")
		}
	})

	t.Run("Socket returns err on failed connection", func(t *testing.T) {
		close := startTestWSServer(t, "test-token", nil)
		defer close()
		socket := socket{
			ip:      "bad ip",
			port:    "bad port",
			testing: true,
		}
		err := socket.connect()
		if err == nil {
			t.Error("Expected connection to fail")
		}
		if socket.connection != nil {
			t.Error("Socket did not close connection after failure")
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("Socket closes connection", func(t *testing.T) {
		close := startTestWSServer(t, "test-token", nil)
		defer close()
		socket := socket{
			ip:      TEST_IP,
			port:    TEST_PORT,
			testing: true,
			token:   "test-token",
		}
		err := socket.connect()
		if err != nil {
			t.Fatalf("Failed to connect socket: %v", err)
		}
		conn := socket.connection
		socket.close()
		if socket.connection != nil {
			t.Error("Socket did not clear connection field")
		}
		err = conn.Close()
		if err == nil {
			t.Error("Socket did not close connection")
		}
	})
}
