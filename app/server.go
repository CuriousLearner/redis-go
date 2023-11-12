package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func parseCommand(encodedCommand []byte) (command string, err error) {
	// TODO: Add validation for correct RESP format

	commandString := string(encodedCommand)
	if commandString[0] != '*' {
		return "", errors.New("Invalid command")
	}
	args := strings.Split(commandString, "\r\n")
	for i, arg := range args[2:] {
		if i%2 == 0 {
			command += arg + " "
		}
	}
	return command, nil
}

func handleConnection(conn net.Conn, buffer []byte) {
	for {
		n, err := conn.Read(buffer)
		if errors.Is(err, io.EOF) {
			fmt.Println("[INFO] EOF detected: ", err.Error())
			break
		}
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
		}
		respEncodedCommand := buffer[:n]
		command, err := parseCommand(respEncodedCommand)
		if err != nil {
			fmt.Println("Error parsing command: ", err.Error())
			break
		}

		response := processCommand(command)
		generateResponse(conn, response)
	}
}

func generateResponse(conn net.Conn, response string) {
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
	}
}

func processCommand(commandString string) (response string) {
	parsedCommand := strings.Split(commandString, " ")
	command := strings.ToUpper(parsedCommand[0])
	args := strings.Trim(strings.Join(parsedCommand[1:], " "), " ")
	switch command {
	case "PING":
		response = "+PONG\r\n"
	case "ECHO":
		response = fmt.Sprintf("$%d\r\n%s\r\n", len(args), args)
	default:
		response = "-ERR unknown command\r\n"
	}
	return response
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

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
