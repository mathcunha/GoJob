package gojob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Slave struct {
	MasterAddr string
	Addr       string
	Name       string
}

func (s *Slave) Work() {
	s.announce()
	http.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		a_path := strings.Split(r.URL.Path, "/")
		if "slave" != a_path[3] {
			http.Error(w, "no handler to path "+r.URL.Path, http.StatusNotFound)
			return
		}
	})
	http.ListenAndServe(s.Addr, nil)
}

//announce itself to master
func (s *Slave) announce() {
	fmt.Printf("hello! I've just arrived master %v\n", s.MasterAddr)
	var postData []byte

	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(s)
	url := fmt.Sprintf("http://%v/api/v1/slave", s.MasterAddr)

	resp, err := http.Post(url, "application/json", w)
	for err != nil {
		fmt.Printf("waiting master at %v\n", s.Addr)
		time.Sleep(1 * time.Minute)
		resp, err = http.Post(url, "application/json", w)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(s)
	if err != nil {
		log.Fatal("ERROR: Master response:", err)

	}

	fmt.Printf("my name is %v\n", s.Name)
}
