package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type macro struct {
	key   string
	delay time.Duration
}

func sendKey(conn *websocket.Conn, key string) error {
	msg := map[string]any{
		"method": "ms.remote.control",
		"params": map[string]any{
			"Cmd":          "Click",
			"DataOfCmd":    key,
			"Option":       "false",
			"TypeOfRemote": "SendRemoteKey",
		},
	}
	err := conn.WriteJSON(msg)
	return err
}

func performMacro(conn *websocket.Conn, macros []macro) error {
	for _, macro := range macros {
		err := sendKey(conn, macro.key)
		if err != nil {
			return err
		}
		time.Sleep(macro.delay)
	}
	return nil
}

func openBrightness(conn *websocket.Conn) error {
	macros := []macro{
		{
			key:   KEY_MORE,
			delay: 2000 * time.Millisecond,
		},
		{
			key:   KEY_ENTER,
			delay: 1000 * time.Millisecond,
		},
		{
			key:   KEY_ENTER,
			delay: 1000 * time.Millisecond,
		},
	}
	err := performMacro(conn, macros)
	return err
}

func closeBrightness(conn *websocket.Conn) error {
	macros := []macro{
		{
			key:   KEY_RETURN,
			delay: 1300 * time.Millisecond,
		},
		{
			key:   KEY_RETURN,
			delay: 1300 * time.Millisecond,
		},
		{
			key:   KEY_RETURN,
			delay: 1300 * time.Millisecond,
		},
	}
	err := performMacro(conn, macros)
	return err
}

func changeBrightness(socket *socket, change int, key string) error {
	if change <= 0 {
		return fmt.Errorf("adjustment value is less than or equal to zero")
	}

	err := socket.connect()
	if err != nil {
		return err
	}

	err = openBrightness(socket.connection)
	if err != nil {
		return err
	}
	for i := 0; i < change; i++ {
		err := sendKey(socket.connection, key)
		if err != nil {
			return err
		}
		time.Sleep(450 * time.Millisecond)
	}
	err = closeBrightness(socket.connection)
	if err != nil {
		return err
	}
	return nil
}
