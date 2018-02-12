package main

import (
	"log"
	"os"
	"yulong-hids/agent/client"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}
	var agent client.Agent
	agent.ServerNetLoc = os.Args[1]
	if len(os.Args) == 3 && os.Args[2] == "debug" {
		log.Println("DEBUG MODE")
		agent.IsDebug = true
	}
	agent.Run()
}
