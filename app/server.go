package main

import (
	"fmt"
	"os"
)

var DataBase = make(map[string]string)

func main() {
	master := make([]string, 2)
	fmt.Println("Logs from your program will appear here!")
	port := "6379"
	options := parseArgs(os.Args[1:])
	if portVal, ok := options["--port"]; ok {
		port = portVal[0]
	}
	if _, ok := options["--replicaof"]; ok {
		master = options["replicaof"]
	}
	protocol := "tcp"
	ipAddr := "0.0.0.0"
	l := listenConnection(protocol, ipAddr, port)
	for {
		conn := acceptConnection(l)
		go respond(conn, master)
	}
}
