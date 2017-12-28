package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "piswitch/statik"

	"github.com/hdra/cron"
	"github.com/rakyll/statik/fs"
	rpio "github.com/stianeikeland/go-rpio"
)

type Schedule struct {
	Id      string   `json:"id"`
	Time    int      `json:"time"`
	Days    []string `json:"days"`
	Command string   `json:"command"`
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
	state   *State
	data    string
	cron    *cron.Cron
	Entries map[CronID]Schedule
}

func (c *Scheduler) DumpData() {
	entries := make([]Schedule, len(scheduler.Entries))
	index := 0
	for _, val := range scheduler.Entries {
		entries[index] = val
		index++
	}
	data, err := json.Marshal(entries)
	if err != nil {
		log.Println("Error serializing data", err.Error())
		return
	}
	err = ioutil.WriteFile(c.data, data, 0644)
	if err != nil {
		log.Println("Error dumping data")
		log.Println(err.Error())
	}
}

func (c *Scheduler) Start() {
	c.cron.Start()
}

func (c *Scheduler) AddSchedule(schedule Schedule) {
	c.Lock()
	id, err := c.cron.AddFunc(schedule.GetCronSpec(), func() {
		fmt.Printf("Executing %v, setting to: %v\n", schedule.Id, schedule.Command)
		switch schedule.Command {
		case "On":
			c.state.On()
		case "Off":
			c.state.Off()
		case "Toggle":
			c.state.Toggle()
		}
		fmt.Println("=================")
	})
	if err != nil {
		panic(err.Error())
	}
	c.Entries[id] = schedule
	c.DumpData()
	c.Unlock()
}

func (c *Scheduler) RemoveSchedule(id string) error {
	c.Lock()
	defer c.Unlock()
	for cronid, schedule := range c.Entries {
		if schedule.Id == id {
			c.cron.Remove(cronid)
			delete(c.Entries, cronid)
			c.DumpData()
			return nil
		}
	}

	return errors.New("Schedule not found")
}

type State struct {
	sync.Mutex
	pin          rpio.Pin
	CurrentState bool
}

func (s *State) Init(pin int) {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.pin = rpio.Pin(pin)
	s.pin.Output()
	s.pin.Low()
}

func (s *State) Cleanup() {
	rpio.Close()
}

func (s *State) Toggle() {
	s.Lock()
	s.pin.Toggle()
	s.CurrentState = !state.CurrentState
	s.Unlock()
}

func (s *State) On() {
	s.Lock()
	s.pin.High()
	s.CurrentState = true
	s.Unlock()
}

func (s *State) Off() {
	s.Lock()
	s.pin.Low()
	s.CurrentState = false
	s.Unlock()
}

func initState() State {
	return State{sync.Mutex{}, 0, false}
}

func initScheduler(data string, state *State) Scheduler {
	jobs := make(map[CronID]Schedule)
	scheduler := Scheduler{
		sync.Mutex{},
		state,
		data,
		cron.New(),
		jobs,
	}

	//load schedules from json
	var schedules []Schedule
	file, err := ioutil.ReadFile(data)
	if err == nil {
		e := json.Unmarshal(file, &schedules)
		if e != nil {
			log.Println(e.Error())
		} else {
			for _, schedule := range schedules {
				scheduler.AddSchedule(schedule)
			}
		}
	}

	scheduler.Start()
	return scheduler
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid method", 400)
		return
	}

	entries := make([]Schedule, len(scheduler.Entries))
	index := 0
	for _, val := range scheduler.Entries {
		entries[index] = val
		index++
	}

	data := struct {
		CurrentState bool       `json:"currentState"`
		Schedules    []Schedule `json:"schedules"`
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
	scheduler.AddSchedule(schedule)
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
	err = scheduler.RemoveSchedule(schedule.Id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

const PIN = 14
const DATAFILE = "data.json"

var scheduler Scheduler
var state State

func main() {
	state = initState()
	state.Init(PIN)
	defer state.Cleanup()
	scheduler = initScheduler(DATAFILE, &state)

	//Get state API
	http.HandleFunc("/api", getIndex)
	//Toggle state API
	http.HandleFunc("/api/toggle", toggleState)
	//Add schedule API
	http.HandleFunc("/api/add", addSchedule)
	//Remove schedule API
	http.HandleFunc("/api/remove", removeSchedule)

	//Client App
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(statikFS))

	http.ListenAndServe(":8000", nil)
}
