package gojob

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Master struct {
	Name    string
	Members []Member
}

//show all announced members
func members() {
}

func (m *Master) addSlave(s *Slave) {
	//TODO need a mutex
	s.Name = fmt.Sprintf("slave_%d", len(Members))
	m.Members = append(m.Members, s)

}

func (m *Master) StartServer() {
	http.HandleFunc("/api/v1/", m.masterRestHandlerV1)
	http.ListenAndServe(getPort(), nil)
}

func (m *Master) masterRestHandlerV1(w http.ResponseWriter, r *http.Request) {
	a_path := strings.Split(r.URL.Path, "/")
	if "slave" == a_path[3] {
		var slave Slave
		if "POST" == r.Method {
			err := json.NewDecoder(r.Body).Decode(&slave)
			if err != nil {
				log.Printf("ERROR: Parsing request body masterRestHandlerV1:%v", err)
				http.Error(w, "ERROR: Parsing request body masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)
			}
			m.addSlave(slave)
			err = json.NewEncoder(w.Body).Encode(slave)
			if err != nil {
				log.Printf("ERROR: Encoding slave masterRestHandlerV1:%v", err)
				http.Error(w, "ERROR: Encoding slave masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)

			}
			w.Header().Set("Content-Type", "application/json")
		}
		if "GET" == r.Method {
			err = json.NewEncoder(w.Body).Encode(m.Members)
			if err != nil {
				log.Printf("ERROR: Encoding slaveS masterRestHandlerV1:%v", err)
				http.Error(w, "ERROR: Encoding slaveS masterRestHandlerV1 "+r.URL.Path, http.StatusInternalServerError)

			}
			w.Header().Set("Content-Type", "application/json")
		}

	}
	http.Error(w, "no handler to path "+r.URL.Path, http.StatusNotFound)
}
