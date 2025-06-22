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

func getOpenMacros(brightnessPosition int, initialDelay time.Duration) []macro {
	openMacros := []macro{{
		key:   KEY_MORE,
		delay: initialDelay,
	}, {
		key:   KEY_ENTER,
		delay: 1000 * time.Millisecond,
	}}
	for i := 0; i < brightnessPosition-1; i++ {
		openMacros = append(openMacros, macro{
			key:   KEY_RIGHT,
			delay: 500 * time.Millisecond,
		})
	}
	openMacros = append(openMacros, macro{
		key:   KEY_ENTER,
		delay: 1000 * time.Millisecond,
	})
	return openMacros
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
		delay := macro.delay
		if isTesting {
			delay = 1 * time.Millisecond
		}
		time.Sleep(delay)
	}
	return nil
}

func changeBrightness(socket *socket, change int, brightnessPosition int, initialDelay time.Duration, key string) error {
	err := socket.connect()
	if err != nil {
		return err
	}

	openMacros := getOpenMacros(brightnessPosition, initialDelay)
	err = performMacro(socket.connection, openMacros)
	if err != nil {
		return err
	}
	changeMacros := []macro{}
	for i := 0; i < change; i++ {
		changeMacros = append(changeMacros, macro{
			key:   key,
			delay: 450 * time.Millisecond,
		})
	}
	err = performMacro(socket.connection, changeMacros)
	if err != nil {
		return err
	}
	err = performMacro(socket.connection, CLOSE_MACRO)
	if err != nil {
		return err
	}
	return nil
}
