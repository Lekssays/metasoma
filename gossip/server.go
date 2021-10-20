package main

import (
	"log"
	"net"
)

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("[*] SERVER: Well received! "), addr)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	address := "0.0.0.0"
	port := 1337
	log.Printf("Listening on %s:%d\n", address, port)
	buffer := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		size, remoteaddr, err := ser.ReadFromUDP(buffer)
		memory := Decode(buffer)
		log.Printf("Read a message of size %d bytes from %v %s \n", size, remoteaddr, memory)
		if err != nil {
			log.Println(err)
			continue
		}
		go sendResponse(ser, remoteaddr)
	}
}
