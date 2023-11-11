package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn, buffer []byte) {
	var err error
	for {
		_, err = conn.Read(buffer)
		if errors.Is(err, io.EOF) {
			fmt.Println("Error reading -- EOF reached: ", err.Error())
			break
		}
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
		}
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	buffer := make([]byte, 1024)

	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, buffer)

	}

}
