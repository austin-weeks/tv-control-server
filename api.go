package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type api struct {
	ws *websocket.Conn
	pw string
}

func getChange(r *http.Request) (int, error) {
	adj := r.Header.Get("Adjustment")
	if adj == "" {
		return 0, fmt.Errorf("no adjustment header")
	}
	change, err := strconv.Atoi(strings.TrimSpace(adj))
	if err != nil {
		return 0, err
	}
	return change, nil
}

func (a *api) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	} else if auth != a.pw {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

func (a *api) increaseBrightness(w http.ResponseWriter, r *http.Request) {
	if ok := a.checkAuth(w, r); !ok {
		return
	}
	change, err := getChange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = changeBrightness(a.ws, change, KEY_RIGHT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *api) decreaseBrightness(w http.ResponseWriter, r *http.Request) {
	if ok := a.checkAuth(w, r); !ok {
		return
	}
	change, err := getChange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = changeBrightness(a.ws, change, KEY_LEFT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
