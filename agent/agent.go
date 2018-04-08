package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"yulong-hids/agent/client"
)

var brokers = flag.String("b", "127.0.0.1:9092", "Kafka brokers")

func main() {
	flag.Parse()
	if len(os.Args) <= 1 || *brokers == "" {
		usage()
		return
	}
	var agent client.Agent
	agent.ServerNetLoc = os.Args[1]
	agent.Brokers = brokers
	if len(os.Args) == 3 && os.Args[2] == "debug" {
		log.Println("DEBUG MODE")
		agent.IsDebug = true
	}
	agent.Run()
}

func usage() {
	fmt.Println("Usage: agent[.exe] ServerIP [debug] -b 127.0.0.1:9092")
	fmt.Println("Example: agent 8.8.8.8 debug")
}
