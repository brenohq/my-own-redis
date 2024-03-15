package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAX_BYTE_SIZE = 512
const DEFAULT_PORT = 6379

type Value struct {
	value        string
	expiry_after int64
	created_at   int64
}

var storage = make(map[string]Value)
var port int

func main() {
	fmt.Println("Logs from your program will appear here!")

	flag.IntVar(&port, "port", DEFAULT_PORT, "Flag used to set which port of this redis instance will run.")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		os.Exit(1)
	}

	defer l.Close()

	for {
		connection, err := l.Accept()

		if err != nil {
			fmt.Print("Error accepting connection: ", err.Error())
			continue
		}

		go handleRequest(connection)
	}
}

func handleRequest(connection net.Conn) {
	defer connection.Close()

	tmp := make([]byte, MAX_BYTE_SIZE)

	for {
		lenght, err := connection.Read(tmp)

		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		input := string(tmp[:lenght])
		inputs := strings.Split(input, "\r\n")
		command, rest := inputs[2], inputs[3:]

		switch lower_command := strings.ToLower(command); lower_command {
		case "ping":
			connection.Write([]byte("+PONG\r\n"))

		case "echo":
			echo := rest[1]
			fmt.Fprintf(connection, "$%d\r\n%s\r\n", len(echo), echo)

		case "set":
			fmt.Println("SET")
			key, value := rest[1], rest[3]
			temp_value := Value{value: value}

			fmt.Println(value)
			fmt.Println(rest)

			if len(rest) >= 6 {
				if expiry_after, err := strconv.Atoi(rest[7]); err == nil {
					fmt.Println(expiry_after)
					temp_value.expiry_after = int64(expiry_after)
				}
			} else {
				temp_value.expiry_after = -1
			}

			connection.Write([]byte("+OK\r\n"))
			now := time.Now()
			temp_value.created_at = now.UnixMilli()
			storage[(key)] = temp_value

		case "get":
			now := time.Now()
			fmt.Println("GET")
			key := rest[1]
			value := storage[(key)]

			fmt.Print(value.expiry_after)
			fmt.Println((now.UnixMilli() - value.created_at))

			if value.expiry_after < 0 || ((now.UnixMilli() - value.created_at) < value.expiry_after) {
				fmt.Fprintf(connection, "$%d\r\n%s\r\n", len(value.value), value.value)
			} else {
				fmt.Fprint(connection, "$-1\r\n")
			}

		case "info":
			var section string

			if len(rest) > 1 {
				section = rest[1]
			} else {
				section = "replication"
			}

			if section == "replication" {
				response := "role:master"
				connection.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(response), response)))
			}
		}
	}
}
