package main

import (
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
			delay: 1500 * time.Millisecond,
		},
		{
			key:   KEY_ENTER,
			delay: 600 * time.Millisecond,
		},
		{
			key:   KEY_ENTER,
			delay: 600 * time.Millisecond,
		},
	}
	err := performMacro(conn, macros)
	return err
}

func closeBrightness(conn *websocket.Conn) error {
	macros := []macro{
		{
			key:   KEY_RETURN,
			delay: 900 * time.Millisecond,
		},
		{
			key:   KEY_RETURN,
			delay: 900 * time.Millisecond,
		},
		{
			key:   KEY_RETURN,
			delay: 900 * time.Millisecond,
		},
	}
	err := performMacro(conn, macros)
	return err
}

func changeBrightness(conn *websocket.Conn, change int, key string) error {
	err := openBrightness(conn)
	if err != nil {
		return err
	}
	for i := 0; i < change; i++ {
		err := sendKey(conn, key)
		if err != nil {
			return err
		}
		time.Sleep(300 * time.Millisecond)
	}
	err = closeBrightness(conn)
	return err
}
