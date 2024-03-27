package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var DataBase = make(map[string]string)

func acceptConnection(listner net.Listener) net.Conn {
	conn, err := listner.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	return conn
}

func respond(conn net.Conn) {
	defer conn.Close()
	for {
		packet := make([]byte, 256)
		var response string
		_, readErr := conn.Read(packet)
		if readErr != nil {
			fmt.Println(readErr.Error())
			return
		}
		elements := redisProtocolParser(packet)
		switch strings.ToLower(elements[2]) {
		case "ping":
			response = string("+PONG\r\n")
		case "echo":
			response = elements[3] + "\r\n" + elements[4] + "\r\n"
		case "set":
			key := elements[4]
			value := elements[5] + "\r\n" + elements[6] + "\r\n"
			for i := 7; i < len(elements)-1; i++ {
				if strings.ToLower(elements[i+1]) == "px" {
					lifetime, err := strconv.Atoi(elements[i+3])
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					go expiry(lifetime, key)
					break
				}
				value += elements[i] + "\r\n"
			}
			DataBase[key] = value
			response = "+OK\r\n"
		case "get":
			key := elements[4]
			response = DataBase[key]
		}
		_, writeErr := conn.Write([]byte(response))
		if writeErr != nil {
			fmt.Println(writeErr.Error())
			return
		}
	}
}

func expiry(lifetime int, key string) {
	time.Sleep(time.Millisecond * time.Duration(lifetime))
	DataBase[key] = "$-1\r\n"
	fmt.Println(DataBase[key])
}
func redisProtocolParser(packet []byte) []string {
	content := packet[2:]
	elements := strings.Split(string(content), "\r\n")
	return elements
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn := acceptConnection(l)
		go respond(conn)
	}
}
