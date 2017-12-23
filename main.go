package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Schedule struct {
	Id      string
	Time    int
	Days    []string
	Command string
}

type State struct {
	sync.Mutex
	CurrentState bool
	Schedules    []Schedule
}

func (s *State) Toggle() {
	s.Lock()
	s.CurrentState = !state.CurrentState
	//Set GPIO pin
	s.Unlock()
}

func getInitialState() State {
	//load json
	//setup pins
	fmt.Println("Getting initial state")
	return State{sync.Mutex{}, false, []Schedule{
		{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
		{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
	}}
}

var state = getInitialState()

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		badMethod(w)
	} else {
		b, _ := json.Marshal(state)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(b))
	}
}

func toggleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		badMethod(w)
	} else {
		state.Toggle()
		w.WriteHeader(http.StatusOK)
	}
}

func badMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "Bad request")
}

func main() {
	//Get state API
	http.HandleFunc("/", getIndex)
	//Toggle state API
	http.HandleFunc("/toggle", toggleState)
	//Add schedule API
	//Remove schedule API
	http.ListenAndServe(":8000", nil)
}
