package main

import (
	"bytes"
	"encoding/json"
	"github.com/mathcunha/GoJob/gojob"
	"io/ioutil"
	"log"
	"os"
)

var profile string
var addr string
var masterAddr string
var configFile string

func init() {
	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = "localhost:8180"
	}
	masterAddr = os.Getenv("MASTER_ADDR")
	if masterAddr == "" {
		masterAddr = "localhost:8080"
	}
	profile = os.Getenv("PROFILE")
	if profile == "" {
		profile = "slave"
	}
	configFile = os.Getenv("CONFIG")
	if configFile == "" {

		configFile = "mongodb/master.json"
	}
}

func loadBenchmark(cfgPath string) (benchmark gojob.Benchmark) {
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

func main() {
	var server gojob.Node
	switch profile {
	case "slave":
		server = &gojob.Slave{Addr: addr, MasterAddr: masterAddr}
	case "master":
		server = gojob.NewMaster(masterAddr, loadBenchmark(configFile))
	}
	server.Start()
}
