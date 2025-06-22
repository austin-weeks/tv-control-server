package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type api struct {
	socket             *socket
	pw                 string
	brightnessPosition int
}

const MAX_CHANGE = 50

func getChange(r *http.Request) (int, error) {
	adj := r.Header.Get("Adjustment")
	if adj == "" {
		return 0, fmt.Errorf("no adjustment header")
	}
	change, err := strconv.Atoi(strings.TrimSpace(adj))
	if err != nil {
		return 0, err
	}
	if change <= 0 || change > 50 {
		return 0, fmt.Errorf("adjusment value is not in range 1 to 50")
	}
	return change, nil
}

func (a *api) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if a.pw == "" {
		return true
	}
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
	err = changeBrightness(a.socket, change, a.brightnessPosition, KEY_RIGHT)
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
	err = changeBrightness(a.socket, change, a.brightnessPosition, KEY_LEFT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
