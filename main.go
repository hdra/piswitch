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

type Scheduler struct {
	sync.Mutex
	Entries []Schedule
}

func (c *Scheduler) AddSchedule(schedule Schedule) {
	c.Lock()
	c.Entries = append(c.Entries, schedule)
	//Update crons
	c.Unlock()
}

func (c *Scheduler) RemoveSchedule(id string) {
	c.Lock()
	for i, schedule := range c.Entries {
		if schedule.Id == id {
			c.Entries = append(c.Entries[:i], c.Entries[i+1:]...)
			break
		}
	}
	//Update crons
	c.Unlock()
}

type State struct {
	sync.Mutex
	CurrentState bool
}

func (s *State) Toggle() {
	s.Lock()
	//Set GPIO pin
	s.CurrentState = !state.CurrentState
	s.Unlock()
}

func (s *State) On() {
	s.Lock()
	// Set pin to HIGH
	s.CurrentState = true
	s.Unlock()
}

func (s *State) Off() {
	s.Lock()
	//Set GPIO pin to LOW
	s.CurrentState = false
	s.Unlock()
}

func getInitialState() State {
	//setup pins
	fmt.Println("Getting initial state")
	return State{sync.Mutex{}, false}
}

func loadSchedules() Scheduler {
	//load schedules from json
	return Scheduler{
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
		http.Error(w, "Invalid method", 400)
		return
	}
	data := struct {
		CurrentState bool
		Schedules    []Schedule
	}{state.CurrentState, schedules.Entries}
	b, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(b))
}

func toggleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", 400)
		return
	}
	state.Toggle()
	w.WriteHeader(http.StatusOK)
}

func addSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", 400)
		return
	}
	if r.Body == nil {
		http.Error(w, "Missing param", 400)
		return
	}
	var schedule Schedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	schedules.AddSchedule(schedule)
	w.WriteHeader(http.StatusCreated)
}

func removeSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", 400)
		return
	}
	if r.Body == nil {
		http.Error(w, "Missing param", 400)
		return
	}
	var schedule Schedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	schedules.RemoveSchedule(schedule.Id)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	//Get state API
	http.HandleFunc("/", getIndex)
	//Toggle state API
	http.HandleFunc("/toggle", toggleState)
	//Add schedule API
	http.HandleFunc("/add", addSchedule)
	//Remove schedule API
	http.HandleFunc("/remove", removeSchedule)
	http.ListenAndServe(":8000", nil)
}
