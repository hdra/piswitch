package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Schedule struct {
	Id      string
	Time    int
	Days    []string
	Command string
}

type State struct {
	CurrentState bool
	Schedules    []Schedule
}

func getInitialState() State {
	//load json
	//setup pins
	fmt.Println("Getting initial state")
	return State{false, []Schedule{
		{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
		{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
	}}
}

var state = getInitialState()

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		badMethod(w)
	}
	b, _ := json.Marshal(state)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(b))
}

func toggleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		badMethod(w)
	}
}

func badMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "Bad request")
}

func main() {
	//Get initial state API

	//Toggle state API

	//Add schedule API

	//Remove schedule API

	http.HandleFunc("/", getIndex)
	http.ListenAndServe(":8000", nil)
}
