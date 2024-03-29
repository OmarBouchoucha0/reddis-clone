package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func readPacket(conn net.Conn, packet []byte) []byte {
	n, readErr := conn.Read(packet)
	if readErr != nil {
		fmt.Println(readErr.Error())
		os.Exit(1)
	}
	return packet[:n]
}

func writePacket(conn net.Conn, response string) {
	_, writeErr := conn.Write([]byte(response))
	if writeErr != nil {
		fmt.Println(writeErr.Error())
		os.Exit(1)
	}
}

func pingResponse() string {
	response := string("+PONG\r\n")
	return response
}

func echoResponse(payload []string) string {
	length := payload[2]
	value := payload[3]
	response := length + "\r\n" + value + "\r\n"
	return response
}

func setResponse(payload []string) string {
	key := payload[3]
	value := payload[4] + "\r\n" + payload[5] + "\r\n"
	if isExpiry(payload) {
		lifetime := parseLifetime(payload)
		fmt.Println("sleeping ...")
		go expiry(lifetime, key)
	}
	DataBase[key] = value
	response := "+OK\r\n"
	return response
}

func isExpiry(payload []string) bool {
	if len(payload) > 9 {
		return strings.ToLower(payload[7]) == "px"
	}
	return false
}

func parseLifetime(payload []string) int {
	lifetime, err := strconv.Atoi(payload[9])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return lifetime
}

func expiry(lifetime int, key string) {
	time.Sleep(time.Millisecond * time.Duration(lifetime))
	DataBase[key] = "$-1\r\n"
}

func getResponse(payload []string) string {
	key := payload[3]
	response := DataBase[key]
	return response
}

func infoResponse(master []string) string {
	var response string
	if master[0] == "" {
		response = "$11\r\nrole:master\r\n"
	} else {
		response = "$10\r\nrole:slave\r\n"
	}
	return response
}

func respond(conn net.Conn, master []string) {
	defer conn.Close()
	packet := make([]byte, 256)
	var response string
	for {
		packet = readPacket(conn, packet)
		payload := parserRedisProtocol(packet)
		command := parseCommandformPayload(payload)
		switch command {
		case "ping":
			response = pingResponse()
		case "echo":
			response = echoResponse(payload)
		case "set":
			response = setResponse(payload)
		case "get":
			response = getResponse(payload)
		case "info":
			response = infoResponse(master)
		default:
			response = "+ERR unknown command\r\n"
		}
		writePacket(conn, response)
	}
}
