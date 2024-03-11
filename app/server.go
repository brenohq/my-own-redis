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

	for {
		connection, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(connection)
	}
}

func handleRequest(connection net.Conn) {
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

		switch strings.ToLower(command[2]) {
		case "echo":
			connection.Write([]byte("$" + strconv.Itoa(len(command[4])) + "\r\n" + string(strings.Join(command[4:], "\r\n"))))
		default:
			connection.Write([]byte("+PONG\r\n"))
		}
	}
}
