package main

import (
	"github.com/mathcunha/GoJob/dummy"
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

		configFile = "dummy/master.json"
	}
}

func loadTask(task string) gojob.Task {
	switch task {
	case "dummy":
		return &dummy.Dummy{}
	}
	return nil
}

func main() {
	var server gojob.Node
	switch profile {
	case "slave":
		slave := gojob.NewSlave(addr, masterAddr)
		slave.Task = loadTask(slave.Benchmark.Workload)
		server = slave
	case "master":
		server = gojob.NewMaster(masterAddr, configFile)
	}
	server.Start()
}
