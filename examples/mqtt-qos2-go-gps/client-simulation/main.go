package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5247")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()
	sendData(conn)
}

func sendData(conn net.Conn) {
	data := []byte("Hello, there!")
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
