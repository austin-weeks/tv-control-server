package main

import (
	"testing"
	"time"
)

func TestPerformMacro(t *testing.T) {
	t.Run("Macro sent correctly", func(t *testing.T) {
		macros := []macro{
			{
				key:   KEY_0,
				delay: 1 * time.Millisecond,
			},
			{
				key:   KEY_RETURN,
				delay: 1 * time.Millisecond,
			},
			{
				key:   KEY_YELLOW,
				delay: 1 * time.Millisecond,
			},
			{
				key:   KEY_OPEN,
				delay: 1 * time.Millisecond,
			},
		}
		close := startTestWSServer(t, "", macros)
		defer close()
		socket, err := getTestSocket()
		if err != nil {
			t.Fatal(err)
		}
		err = performMacro(socket.connection, macros)
		if err != nil {
			t.Error("Error performing macros:", err)
		}
	})
}

func TestChangeBrightnes(t *testing.T) {
	t.Run("Macros constructed correctly", func(t *testing.T) {
		changeMacros := []macro{
			{
				key:   KEY_RIGHT,
				delay: 1 * time.Millisecond,
			},
			{
				key:   KEY_RIGHT,
				delay: 1 * time.Millisecond,
			},
			{
				key:   KEY_RIGHT,
				delay: 1 * time.Millisecond,
			},
		}
		expectedMacros := append(append(OPEN_MACRO, changeMacros...), CLOSE_MACRO...)
		close := startTestWSServer(t, "", expectedMacros)
		defer close()
		socket, err := getTestSocket()
		if err != nil {
			t.Fatal(err)
		}
		err = changeBrightness(&socket, 3, KEY_RIGHT)
		if err != nil {
			t.Error("Error changing brightness:", err)
		}
	})
}
