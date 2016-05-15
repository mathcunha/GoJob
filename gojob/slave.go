package gojob

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Slave struct {
	MasterEndPoint endPoint
	EndPoint       endPoint
	Name           string
}

func (s *Slave) StartServer() {
	http.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		a_path := strings.Split(r.URL.Path, "/")
		if "slave" != a_path[3] {
			http.Error(w, "no handler to path "+r.URL.Path, http.StatusNotFound)
			return
		}
	})
	http.ListenAndServe(getPort(), nil)
}

func NewSlave() Slave {
	return Slave{name: os.Hostname()}
}

//announce itself to master
func (s *Slave) announce() {
	fmt.Printf("hello! I've just arrived master %v\n", s.MasterEndPoint)
	var postData []byte

	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(s)
	url := fmt.Sprintf("http://%v:%d/alert", s.MasterEndPoint.Host, s.MasterEndPoint.Port)

	resp, err := http.Post(url, "application/json", w)
	for err != nil {
		fmt.Printf("waiting master at %v\n", s.Endpoint)
		time.Sleep(1 * time.Minute)
		resp, err = http.Post(url, "application/json", w)
	}
	defer resp.Body.Close()

	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		log.Printf("ERROR: Parsing request body announce:%v", err)
		http.Error(w, "ERROR: Parsing request body announce"+r.URL.Path, http.StatusInternalServerError)
	}

	fmt.Printf("my name is %v\n", s.Name)
}
