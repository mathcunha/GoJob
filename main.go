package main

import (
	"github.com/mathcunha/GoJob/gojob"
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

		configFile = "config.json"
	}
}

func main() {
	var server gojob.Worker
	switch profile {
	case "slave":
		server = &gojob.Slave{Addr: addr, MasterAddr: masterAddr}
	case "master":
		server = &gojob.Master{Name: "master", Addr: masterAddr}
	}
	server.Work()
}
