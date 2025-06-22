package main

import (
	"testing"
	"time"
)

func TestGetOpenMacro(t *testing.T) {
	t.Run("Correct with location <= 1", func(t *testing.T) {
		for _, loc := range []int{-1, 0, 1} {
			openMacros := getOpenMacros(loc, 1)
			if len(openMacros) != 3 {
				t.Error("Expected macros to have length of 3:", openMacros)
			}
		}
	})

	t.Run("Correct with location > 1", func(t *testing.T) {
		for _, loc := range []int{2, 3, 4} {
			openMacros := getOpenMacros(loc, 1)
			expectedLen := 3 + (loc - 1)
			if len(openMacros) != expectedLen {
				t.Errorf("Expected macros to have length of %d: %v", expectedLen, openMacros)
			}
		}
	})

	t.Run("Correct initial delay", func(t *testing.T) {
		initialDelay := 5*time.Second + 3*time.Millisecond
		openMacros := getOpenMacros(1, initialDelay)
		if openMacros[0].delay != initialDelay {
			t.Errorf("Expected delay of %d != actual delay of %d", initialDelay, openMacros[0].delay)
		}
	})
}

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
		brightLoc := 1
		expectedMacros := append(append(getOpenMacros(brightLoc, 1), changeMacros...), CLOSE_MACRO...)
		close := startTestWSServer(t, "", expectedMacros)
		defer close()
		socket, err := getTestSocket()
		if err != nil {
			t.Fatal(err)
		}
		err = changeBrightness(&socket, 3, brightLoc, 1, KEY_RIGHT)
		if err != nil {
			t.Error("Error changing brightness:", err)
		}
	})
}
