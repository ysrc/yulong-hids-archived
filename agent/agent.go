package main

import (
	"fmt"
	"log"
	"os"
	"yulong-hids/agent/client"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: agent[.exe] ServerIP [debug]")
		fmt.Println("Example: agent 8.8.8.8 debug")
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
