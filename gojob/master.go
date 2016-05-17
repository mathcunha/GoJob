package gojob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Master struct {
	//Master name
	Name string
	//Slaves registered
	Members []Slave
	//Masters http server address
	Addr      string
	mutex     *sync.Mutex
	Benchmark Benchmark
}

func loadBenchmark(cfgPath string) (benchmark Benchmark) {
	file, err := os.Open(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&benchmark)
	return
}

type Benchmark struct {
	//Number of nodes running the TestCase
	Nodes int
	//Users (goroutines) per node
	Users int
	//Operation per user
	Ops int
	//TODO work in progress
	Workload string
	//Some variables needed by Tasks
	Properties map[string]string
}

type Node interface {
	Start()
}

//The right way to create a new master
func NewMaster(addr string, cfgFile string) *Master {
	return &Master{Name: "master", Addr: addr, mutex: &sync.Mutex{}, Benchmark: loadBenchmark(cfgFile)}
}

func (m *Master) addSlave(s *Slave) (added bool) {
	s.Name = ""
	m.mutex.Lock()
	members := len(m.Members)
	if m.Benchmark.Nodes > members {
		s.Name = fmt.Sprintf("slave_%d", len(m.Members))
		m.Members = append(m.Members, *s)
		added = true
	}
	m.mutex.Unlock()
	if added {
		s.Benchmark = m.Benchmark
	}
	return added
}

func (m *Master) Start() {
	http.HandleFunc("/slave", m.slaveHandlerV1)
	http.ListenAndServe(m.Addr, nil)
}

func (m *Master) NotifySlaves() {
	var wg sync.WaitGroup
	wg.Add(len(m.Members))
	var postData []byte
	w := bytes.NewBuffer(postData)
	for _, slave := range m.Members {
		go func(s Slave) {
			url := fmt.Sprintf("http://%v/benchmark?start=true", s.Addr)

			resp, err := http.Post(url, "application/json", w)
			if err != nil {
				log.Printf("error starting slave %v at %v\n", s.Name, s.Addr)
			} else {
				defer resp.Body.Close()
			}
			wg.Done()
		}(slave)
	}
	wg.Wait()
}

func (m *Master) slaveHandlerV1(w http.ResponseWriter, r *http.Request) {
	var slave Slave
	if "POST" == r.Method {
		err := json.NewDecoder(r.Body).Decode(&slave)
		if err != nil {
			log.Printf("ERROR: Parsing request body masterRestHandlerV1:%v", err)
			http.Error(w, "ERROR: Parsing request body masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		if added := m.addSlave(&slave); added {
			err = json.NewEncoder(w).Encode(slave)
			if err != nil {
				log.Printf("ERROR: Encoding slave masterRestHandlerV1:%v", err)
				http.Error(w, "ERROR: Encoding slave masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)
				return

			}
			m.mutex.Lock()
			if len(m.Members) == m.Benchmark.Nodes {
				//send response to slave and wait 10 seconds before notifyAll
				go func() {
					time.Sleep(10 * time.Second)
					m.NotifySlaves()
				}()
			}
			m.mutex.Unlock()
		} else {
			fmt.Fprintf(w, fmt.Sprintf("{\"error\":\"already full of slaves. number of nodes equal to %d\"}", m.Benchmark.Nodes))
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}
	if "GET" == r.Method {
		err := json.NewEncoder(w).Encode(m.Members)
		if err != nil {
			log.Printf("ERROR: Encoding slaveS masterRestHandlerV1:%v", err)
			http.Error(w, "ERROR: Encoding slaveS masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}

	http.Error(w, r.Method, http.StatusMethodNotAllowed)
}
