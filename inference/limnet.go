package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// print(C.initialize(false))
	
	address := "0.0.0.0"
	port := 1337
	log.Printf("Listening on %s:%d\n", address, port)
	buffer := make([]byte, 512)
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
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(buffer, size, remoteaddr)
		var numberOfEntries int = ParseBuffer(buffer, size)
		fmt.Println("Number of Memories =",numberOfEntries)
		// on_memories_received(&buffer[0], &buffer[4*num_entries], num_entries)
	}
}

func BytesToI32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func ParseBuffer(buffer []byte, size int) int {
	var bufferedMemories [10][]byte
	memoriesCount := 0
	var start uint32 = 0
	var end uint32 = 0
	for end < uint32(size) {
		var memorySize uint32
		memorySize = BytesToI32(buffer[start:start+4])
		end = start + uint32(8) + (memorySize * uint32(4))
		bufferedMemories[memoriesCount] = buffer[start:end]
		start = end
		fmt.Println(bufferedMemories[memoriesCount])
		memoriesCount += 1
	}
	return memoriesCount
}

// func send_memories(conn *net.UDPConn) {
// 	peer := net.UDPAddr{
// 		Port: 1337,
// 		IP:  C.get_random_peer(),
// 	}
// 	num_entries = 10
// 	buffer := make([]byte, num_entries*(4+C.compressed_memory_size()))
// 	num_entries = C.get_memories_to_share(&buffer[0], &buffer[4*num_entries], num_entries)
// 	_, err := conn.WriteToUDP(buffer.slice(0, num_entries*(4+C.compressed_memory_size())), peer)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
