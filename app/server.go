package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conns := make(chan net.Conn)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}
			conns <- conn
		}
	}()

	defer l.Close()

	for c := range conns {
		go handleRequest(c)
	}
}

func handleRequest(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		_, err := c.Read(buf)

		if err != nil {
			fmt.Println("Error Reading: ", err.Error())
			os.Exit(1)
		}

		pattern := regexp.MustCompile(`(\*\d+\r?\n\$\d+\r\n)|(\r?\n\$\d+\r\n)|(\r\n)`)
		command := strings.Fields(strings.Join(pattern.Split(string(buf), -1), " "))

		switch command[0] {
		case "ping":
			msg := "+PONG\r\n"
			_, err := c.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
		case "echo":
			_, err := c.Write([]byte("+" + command[1] + "\r\n"))
			if err != nil {
				fmt.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
		}
	}
}
