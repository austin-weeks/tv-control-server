package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	isTesting = true
	code := m.Run()
	os.Exit(code)
}
