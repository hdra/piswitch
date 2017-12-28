package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/hdra/cron"
)

type Schedule struct {
	Id      string
	Time    int
	Days    []string
	Command string
}

func (s Schedule) GetCronSpec() string {
	hour := s.Time / 3600
	minutes := (s.Time - (hour * 3600)) / 60
	seconds := s.Time - (hour * 3600) - (minutes * 60)

	days := make([]string, len(s.Days))
	for i, v := range s.Days {
		days[i] = v[:3]
	}

	return fmt.Sprintf("%v %v %v * * %v", seconds, minutes, hour, strings.Join(days, ","))
}

type CronID = cron.EntryID
type Scheduler struct {
	sync.Mutex
	cron    *cron.Cron
	Entries map[CronID]Schedule
}

func (c *Scheduler) Start() {
	c.cron.Start()
}

func (c *Scheduler) AddSchedule(schedule Schedule) {
	c.Lock()
	id, err := c.cron.AddFunc(schedule.GetCronSpec(), func() {
		fmt.Printf("Running for %v, setting to: %v\n", schedule.Id, schedule.Command)
		switch schedule.Command {
		case "On":
			state.On()
		case "Off":
			state.Off()
		case "Toggle":
			state.Toggle()
		}
		fmt.Println("=================")
	})
	if err != nil {
		panic(err.Error())
	}
	c.Entries[id] = schedule
	c.Unlock()
}

func (c *Scheduler) RemoveSchedule(id string) error {
	c.Lock()
	defer c.Unlock()
	for cronid, schedule := range c.Entries {
		if schedule.Id == id {
			c.cron.Remove(cronid)
			delete(c.Entries, cronid)
			return nil
		}
	}

	return errors.New("Schedule not found")
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
	jobs := make(map[CronID]Schedule)
	scheduler := Scheduler{
		sync.Mutex{},
		cron.New(),
		jobs,
	}

	//load schedules from json
	schedules := []Schedule{
		{"abc", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
		{"def", 1223, []string{"Monday", "Tuesday", "Wednesday"}, "On"},
	}
	for _, schedule := range schedules {
		scheduler.AddSchedule(schedule)
	}

	scheduler.Start()
	return scheduler
}

var state = getInitialState()
var schedules = loadSchedules()

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", 400)
		return
	}

	entries := make([]Schedule, len(schedules.Entries))
	index := 0
	for _, val := range schedules.Entries {
		entries[index] = val
		index++
	}

	data := struct {
		CurrentState bool
		Schedules    []Schedule
	}{state.CurrentState, entries}
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
	err = schedules.RemoveSchedule(schedule.Id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
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
