package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	storage := make(map[string]string)

	for {
		connection, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(connection, storage)
	}
}

func handleRequest(connection net.Conn, storage map[string]string) {
	defer connection.Close()

	for {
		command_byte := make([]byte, 1024)
		command_length, err := connection.Read(command_byte)

		if err != nil {
			fmt.Println("Error reading the request buffer", err)
			return
		}

		command := strings.Split(string(command_byte[:command_length]), "\r\n")

		fmt.Println("---------------------------------")
		fmt.Println(command)
		fmt.Println("---------------------------------")

		var response []byte

		switch strings.ToLower(command[2]) {
		case "echo":
			response = []byte("$" + strconv.Itoa(len(command[4])) + "\r\n" + string(strings.Join(command[4:], "\r\n")))
		case "set":
			storage[command[4]] = command[6]
			response = []byte("+OK\r\n")
		case "get":
			response = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(storage[command[4]]), storage[command[4]]))
		default:
			response = []byte("+PONG\r\n")
		}

		connection.Write(response)
	}
}
