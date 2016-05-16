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
}

func (s *Slave) Start() {
	s.announce()
	http.HandleFunc("/benchmark", func(w http.ResponseWriter, r *http.Request) {
		if "POST" == r.Method {
			start := r.URL.Query().Get("start")
			if "true" == start {
				log.Printf("%v starting benchmark", s.Name)
			}
		} else if "GET" == r.Method {
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
