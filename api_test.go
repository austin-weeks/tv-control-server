package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const DUMMY_URL = "http://www.example.com"

func TestGetChange(t *testing.T) {
	t.Run("Checks for missing adjustment header", func(t *testing.T) {
		req := httptest.NewRequest("", DUMMY_URL, nil)
		_, err := getChange(req)
		if err == nil {
			t.Error("Expected error")
		}
	})

	for _, val := range []string{"not a number", "1.2"} {
		t.Run("Checks for non int adjustment value", func(t *testing.T) {
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Adjustment", val)
			_, err := getChange(req)
			if err == nil {
				t.Error("Expected error")
			}
		})
	}

	for _, val := range []string{"-2", "0", "51"} {
		t.Run("Checks for out of range adjustment value", func(t *testing.T) {
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Adjustment", val)
			_, err := getChange(req)
			if err == nil {
				t.Error("Expected error")
			}
		})
	}

	for _, tc := range []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"3", 3},
		{"5", 5},
	} {
		t.Run("Extracts valid adjustment value", func(t *testing.T) {
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Adjustment", tc.input)
			res, err := getChange(req)
			if err != nil {
				t.Error(err)
			}
			if res != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, res)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	t.Run("Pass if no password required", func(t *testing.T) {
		api := api{
			pw: "",
		}
		req := httptest.NewRequest("", DUMMY_URL, nil)
		w := httptest.NewRecorder()

		res := api.checkAuth(w, req)
		if res != true {
			t.Error("Expected auth to pass")
		}
		if w.Result().StatusCode != http.StatusOK {
			t.Error("Expected status 200")
		}
	})

	t.Run("Pass with matching password", func(t *testing.T) {
		password := "password123"
		api := api{
			pw: password,
		}
		req := httptest.NewRequest("", DUMMY_URL, nil)
		req.Header.Set("Authorization", password)
		w := httptest.NewRecorder()

		res := api.checkAuth(w, req)
		if res != true {
			t.Fail()
		}
		if w.Result().StatusCode != http.StatusOK {
			t.Fail()
		}
	})

	t.Run("Rejects missing authorization header", func(t *testing.T) {
		api := api{
			pw: "password123",
		}
		req := httptest.NewRequest("", DUMMY_URL, nil)
		w := httptest.NewRecorder()

		res := api.checkAuth(w, req)
		if res != false {
			t.Fail()
		}
		if w.Result().StatusCode != http.StatusUnauthorized {
			t.Fail()
		}
	})

	t.Run("Fails with incorrect password", func(t *testing.T) {
		api := api{
			pw: "password123",
		}
		req := httptest.NewRequest("", DUMMY_URL, nil)
		req.Header.Set("Authorization", "incorrect-password")
		w := httptest.NewRecorder()

		res := api.checkAuth(w, req)
		if res != false {
			t.Fail()
		}
		if w.Result().StatusCode != http.StatusForbidden {
			t.Fail()
		}
	})
}

func TestBrightnessHandlers(t *testing.T) {
	close := startTestWSServer(t, "", nil)
	defer close()
	socket, err := getTestSocket()
	if err != nil {
		t.Fatal("Could not create test socket:", err)
	}
	t.Run("Excercise happy path", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			password := "password123"
			api := api{
				pw:     password,
				socket: &socket,
			}
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Authorization", password)
			req.Header.Set("Adjustment", "5")
			w := httptest.NewRecorder()
			if i == 0 {
				api.increaseBrightness(w, req)
			} else {
				api.decreaseBrightness(w, req)
			}
			if w.Result().StatusCode != http.StatusOK {
				t.Error("Expected requests to succeed:", w.Result().StatusCode)
			}
		}
	})

	t.Run("Handlers reject unauthorized requests", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			password := "password123"
			api := api{
				pw:     password,
				socket: &socket,
			}
			req := httptest.NewRequest("", DUMMY_URL, nil)
			w := httptest.NewRecorder()
			if i == 0 {
				api.increaseBrightness(w, req)
			} else {
				api.decreaseBrightness(w, req)
			}
			if w.Result().StatusCode == http.StatusOK {
				t.Error("Expected requests to be rejected")
			}
		}
	})

	t.Run("Handlers reject requests with bad adjustment", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			password := "password123"
			api := api{
				pw:     password,
				socket: &socket,
			}
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Authorization", password)
			w := httptest.NewRecorder()
			if i == 0 {
				api.increaseBrightness(w, req)
			} else {
				api.decreaseBrightness(w, req)
			}
			if w.Result().StatusCode != http.StatusBadRequest {
				t.Error("Expected requests to be rejected")
			}
		}
	})
}

func TestHandlersNoConnection(t *testing.T) {
	t.Run("Endpoints handler failed connection", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			socket := socket{
				ip:   TEST_IP,
				port: TEST_PORT,
			}
			api := api{
				pw:     "password123",
				socket: &socket,
			}
			req := httptest.NewRequest("", DUMMY_URL, nil)
			req.Header.Set("Authorization", "password123")
			req.Header.Set("Adjustment", "5")
			w := httptest.NewRecorder()
			if i == 0 {
				api.increaseBrightness(w, req)
			} else {
				api.decreaseBrightness(w, req)
			}
			if w.Result().StatusCode != http.StatusInternalServerError {
				t.Errorf("Expected requests to be rejected: %v", w.Result().StatusCode)
			}
		}
	})
}
