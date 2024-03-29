package main

import (
	"fmt"
	"net"
	"os"
)

func acceptConnection(listner net.Listener) net.Conn {
	conn, err := listner.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	return conn
}

func listenConnection(protocol string, ipAddr string, port string) net.Listener {
	listner, err := net.Listen(protocol, ipAddr+":"+port)
	if err != nil {
		fmt.Printf("Failed to bind to port %v", port)
		os.Exit(1)
	}
	return listner
}
