package main

import "fmt"

func parseArgs(args []string) map[string][]string {
	options := make(map[string][]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--port":
			if i+1 < len(args) {
				options[arg] = []string{args[i+1]}
				i++
			} else {
				fmt.Println("Missing value for --port option")
			}
		case "--replicaof":
			if i+2 < len(args) {
				options[arg] = []string{args[i+1], args[i+2]}
				i += 2
			} else {
				fmt.Println("Missing values for replicaof option")
			}
		default:
			fmt.Printf("Unknown option: %s\n", arg)
		}
	}
	return options
}
