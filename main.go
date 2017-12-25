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

type SchedulesContainer struct {
	sync.Mutex
	Entries []Schedule
}

type State struct {
	sync.Mutex
	CurrentState bool
}

func (s *State) Toggle() {
	s.Lock()
	s.CurrentState = !state.CurrentState
	//Set GPIO pin
	s.Unlock()
}

func getInitialState() State {
	//setup pins
	fmt.Println("Getting initial state")
	return State{sync.Mutex{}, false}
}

func loadSchedules() SchedulesContainer {
	//load schedules from json
	return SchedulesContainer{
		sync.Mutex{},
		[]Schedule{
			{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
			{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
		}}
}

var state = getInitialState()
var schedules = loadSchedules()

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		badMethod(w)
	} else {
		data := struct {
			CurrentState bool
			Schedules    []Schedule
		}{state.CurrentState, schedules.Entries}
		b, _ := json.Marshal(data)
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
