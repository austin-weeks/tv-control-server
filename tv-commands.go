package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type macro struct {
	key   string
	delay time.Duration
}

type keyMsg struct {
	Method string       `json:"method"`
	Params keyMsgParams `json:"params"`
}
type keyMsgParams struct {
	Cmd          string `json:"Cmd"`
	DataOfCmd    string `json:"DataOfCmd"`
	Option       bool   `json:"Option"`
	TypeOfRemote string `json:"TypeOfRemote"`
}

var OPEN_MACRO = []macro{
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

var CLOSE_MACRO = []macro{
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

func sendKey(conn *websocket.Conn, key string) error {
	msg := keyMsg{
		Method: "ms.remote.control",
		Params: keyMsgParams{
			Cmd:          "Click",
			DataOfCmd:    key,
			Option:       false,
			TypeOfRemote: "SendRemoteKey",
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
	err := performMacro(conn, OPEN_MACRO)
	return err
}

func closeBrightness(conn *websocket.Conn) error {
	err := performMacro(conn, CLOSE_MACRO)
	return err
}

func changeBrightness(socket *socket, change int, key string) error {
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
