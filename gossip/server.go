package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func sendResponse(conn net.Conn) {
	_, err := conn.Write([]byte("[*] SERVER: Well received!"))
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Printf("Handling request from %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 2048)
	for {
		len, err := conn.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		message := string(buffer[:len])

		fmt.Println(message)

		go sendResponse(conn)
	}
}

func main() {
	address := "0.0.0.0"
	port := 1337
	addr := fmt.Sprintf("%s:%d", address, port)
	ser, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ser.Close()

	log.Printf("Listening on %s:%d\n", address, port)
	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Printf("Error accepting connection:", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
