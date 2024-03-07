package main

import (
	"fmt"
	"net"
	"os"
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
	readByte := make([]byte, 100)

	for {
		_, err := c.Read(readByte)
		if err != nil {
			fmt.Println("Error Reading: ", err.Error())
			os.Exit(1)
		}
		c.Write([]byte("+PONG\r\n"))
	}
}
