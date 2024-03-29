package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RespDataType struct {
	value    string
	Datatype string
	category string
}

type redisElement struct {
	value  string
	length uint8
}

var RESPDATATYPELOOKUPMAP = map[string]RespDataType{
	"+": {"+", "Simple String", "Simple"},
	"-": {"-", "Simple Error", "Simple"},
	":": {":", "Integer", "Simple"},
	"$": {"$", "Bulk String", "Aggregate"},
	"*": {"*", "Array", "Aggregate"},
	"_": {"_", "Null", "Simple"},
	"#": {"#", "Boolean", "Simple"},
	",": {",", "Double", "Simple"},
	"(": {"(", "Big Number", "Simple"},
	"!": {"!", "Bulk Error", "Aggregate"},
	"=": {"=", "Verbatim String", "Aggregate"},
	"%": {"%", "Map", "Aggregate"},
	"~": {"~", "Set", "Aggregate"},
	">": {">", "Push", "Aggregate"},
}

func parseElements(packet []byte) []string {
	return strings.Split(string(packet), "\r\n")
}

func extractHeader(elements []string) redisElement {
	dataType := elements[0][0]
	length, err := strconv.Atoi(elements[0][1:])
	if err != nil {
		fmt.Println("coudn extract header")
		os.Exit(1)
	}
	headerElement := redisElement{
		string(dataType),
		uint8(length),
	}
	return headerElement
}

func isAggregate(headerElement redisElement) bool {
	return RESPDATATYPELOOKUPMAP[headerElement.value].category == "Aggregate"
}

func lengthOfElements(elements []string, index uint8) uint8 {
	numberOfElements, err := strconv.Atoi(elements[index][1:])
	if err != nil {
		fmt.Println("cannot parse packet!")
		os.Exit(1)
	}
	return uint8(2 * numberOfElements)
}

func parsePayload(numberOfElements uint8, elements []string) []string {
	var payload []string
	start := uint8(1)
	end := start + uint8(numberOfElements)
	for i := start; i < end; i++ {
		payload = append(payload, elements[i])
	}
	return payload
}

func parseCommandformPayload(payload []string) string {
	command := strings.ToLower(payload[1])
	return command
}

func parseOption(numberOfElements uint8, elements []string) string {
	return "todo"
}

func parserRedisProtocol(packet []byte) []string {
	elements := parseElements(packet)
	headerElement := extractHeader(elements)
	if isAggregate(headerElement) {
		fmt.Println(elements)
		numberOfElements := lengthOfElements(elements, 0)
		fmt.Println(numberOfElements)
		payloadElements := parsePayload(numberOfElements, elements)
		fmt.Println(payloadElements)
		return payloadElements
	} else {
		return []string{elements[0], elements[1], elements[2]}
	}
}
