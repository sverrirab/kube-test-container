package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"log"
	"os"
	"time"
	"text/template"

    "github.com/gorilla/mux"
)

const VERSION = "v1.0"

type State struct {
	template 	*template.Template
	vars 		map[string]interface{}
	abnormal 	bool
	used_ram	[]byte
}


func NewState() *State {
	s := new(State)
	t, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	s.template = t

	s.vars = make(map[string]interface{})
	s.vars["Title"] = "Kube-Test-Container"
	s.vars["Status"] = "Normal load"
	s.vars["Cnt"] = 0

	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	s.vars["Name"] = h
	env := ""
	for _, item := range os.Environ() {
		env += item
	}
	s.vars["Env"] = os.Environ()

	s.used_ram = make([]byte, 0)
	return s
}


func (s *State) Increment(page string) {
	cnt := s.vars["Cnt"].(int)
	s.vars["Cnt"] = cnt + 1
	fmt.Printf("%s - %d\n", page, s.vars["Cnt"] )
}


func (s *State) UseCpu() {
	randomSleep := int64(rand.Intn(4 * 10))
	sleep := time.Duration(randomSleep * 25 * 1000 * 1000)
	time.Sleep(sleep)	// Sleep 0..5 seconds before starting.
	fmt.Println("Started a CPU eating go routine")
	for {
		for i := 0; i < 100000; i++ {
			// empty
		}
		time.Sleep(5 * 1000 * 1000 * 1000)
		go s.UseCpu()
	}
}


func (s *State) UseRam() {
	const MiB = 1024 * 1024
	const size = 10 * MiB
	append_ram := make([]byte, size)
	rand.Read(append_ram)
	s.used_ram = append(s.used_ram, append_ram...)
	s.vars["Status"] = fmt.Sprintf("Abnormal RAM (extra %dMib)", len(s.used_ram) / MiB)
	fmt.Println(s.vars["Status"])
	time.AfterFunc(5 * 1000 * 1000 * 1000, s.UseRam)
}


func (s *State) Index(w http.ResponseWriter, req *http.Request) {
	s.Increment("/")
	s.template.Execute(w, s.vars)
}


func (s *State) Cpu(w http.ResponseWriter, req *http.Request) {
	s.Increment("/cpu")
	if ! s.abnormal {
		s.abnormal = true
		s.vars["Status"] = "Abnormal CPU usage"
		go s.UseCpu()
	}
	s.Status(w, req)
}


func (s *State) Ram(w http.ResponseWriter, req *http.Request) {
	s.Increment("/ram")
	if ! s.abnormal {
		s.abnormal = true
		s.vars["Status"] = "Abnormal RAM"
		go s.UseRam()
	}
	s.Status(w, req)
}


func (s *State) Status(w http.ResponseWriter, req *http.Request) {
	s.Increment("/status")
	fmt.Fprintf(w, "%s\n%s\nrequest #%d %s", s.vars["Name"], s.vars["Status"], s.vars["Cnt"], VERSION)
}


func main() {
	s := NewState()
    router := mux.NewRouter()
    router.HandleFunc("/", s.Index).Methods("GET")
    router.HandleFunc("/cpu", s.Cpu).Methods("GET")
    router.HandleFunc("/ram", s.Ram).Methods("GET")
    router.HandleFunc("/status", s.Status).Methods("GET")

	fmt.Println("Ready")
    log.Fatal(http.ListenAndServe(":8000", router))
}
