package gojob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Slave struct {
	//Master Node address
	MasterAddr string
	//Slave Node address
	Addr string
	//Name given by Master
	Name      string
	Benchmark Benchmark
	workers   []*Worker
	Task      Task
}

func NewSlave(addr, masterAddr string) *Slave {
	slave := Slave{Addr: addr, MasterAddr: masterAddr}
	slave.announce()
	return &slave
}

func (s *Slave) Start() {
	http.HandleFunc("/benchmark", func(w http.ResponseWriter, r *http.Request) {
		if "POST" == r.Method {
			start := r.URL.Query().Get("start")
			if "true" == start {
				log.Printf("%v is starting benchmark for task %T\n", s.Name, s.Task)
				s.Task.LoadProperties(s.Benchmark.Properties)
				for i := 0; i < s.Benchmark.Users; i++ {
					worker := &Worker{Ops: s.Benchmark.Ops, task: s.Task.Clone()}
					s.workers = append(s.workers, worker)
				}
				log.Printf("%d workers loaded\n%v\n", len(s.workers), s.workers)
				for _, worker := range s.workers {
					go func(w *Worker) {
						w.Work()
					}(worker)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, fmt.Sprintf("{\"success\":\"%d workers where started\"}", len(s.workers)))
			}
		} else if "GET" == r.Method {
			result := struct {
				Workers []*Worker
				Name    string
			}{s.workers, s.Name}
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Printf("ERROR: Encoding slave data:%v", err)
				http.Error(w, "ERROR: Encoding slave data ", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, r.Method, http.StatusMethodNotAllowed)
		}
		return
	})
	http.ListenAndServe(s.Addr, nil)
}

//Registering Slave Node
func (s *Slave) announce() {
	log.Printf("hello! I've just arrived master %v\n", s.MasterAddr)
	var postData []byte

	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(s)
	url := fmt.Sprintf("http://%v/slave", s.MasterAddr)

	resp, err := http.Post(url, "application/json", w)
	for err != nil {
		log.Printf("waiting master at %v\n", s.Addr)
		time.Sleep(1 * time.Minute)
		resp, err = http.Post(url, "application/json", w)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		log.Fatal("ERROR: Master response:", err)

	}
	if s.Name == "" {
		log.Fatal("Slave not added to TestScenario which seems full")
	}
	log.Printf("my name is %v\n", s.Name)
}
