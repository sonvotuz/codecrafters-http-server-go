package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// _, err = fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
		// if err != nil {
		// 	fmt.Println("Error writing response:", err.Error())
		// 	os.Exit(1)
		// }

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println("Failed to read from connection", err)
		return
	}

	data := strings.Split(string(req), "\r\n")
	path := strings.Split(data[0], " ")
	if path[1] == "/" {
		response := []byte("HTTP/1.1 200 OK\r\n\r\n")
		conn.Write(response)
	} else {
		response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write(response)
	}
}
