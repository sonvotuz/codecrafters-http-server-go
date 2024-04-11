package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const ListenAddress = "0.0.0.0:4221"

var directoryFlag = flag.String("directory", "", "a string")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ListenAddress)
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
	} else if strings.HasPrefix(path[1], "/echo") {
		randomStr := path[1][6:]

		responseWithContent(conn, randomStr, "text/plain")
	} else if strings.HasPrefix(path[1], "/user-agent") {
		userAgent := strings.Split(data[2], " ")[1]

		responseWithContent(conn, userAgent, "text/plain")
	} else if strings.HasPrefix(path[1], "/files") {
		fileName := path[1][7:]

		dataBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", *directoryFlag, fileName))
		if err != nil {
			response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
			conn.Write(response)
			return
		}
		data := string(dataBytes)
		responseWithContent(conn, data, "application/octet-stream")
	} else {
		response := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write(response)
	}
}

func responseWithContent(conn net.Conn, data, contentType string) {
	buf := bytes.Buffer{}
	buf.WriteString("HTTP/1.1 200 OK\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data)))
	buf.WriteString(fmt.Sprintf("%s", data))

	fmt.Println("Running here")

	conn.Write(buf.Bytes())
}
