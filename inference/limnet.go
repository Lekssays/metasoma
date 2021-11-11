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

type Memory struct {
	Owner	uint32
	Content []uint32
}

func BytesToI32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func ParseBuffer(buffer []byte, size int) (uint32, []byte, []byte) {
	const memElementsSize = 32
	var memoriesCount uint32
	memoriesCount = BytesToI32(buffer[0:4])
	var ownersBytes []byte = buffer[4:4+memoriesCount*4]
	var contentsBytes []byte = buffer[4+memoriesCount*4:4+memoriesCount*4+memoriesCount*memElementsSize*4]
	return memoriesCount, ownersBytes, contentsBytes
}

func ParseMemories(memoriesCount uint32, ownersBytes []byte, memoriesBytes []byte) []Memory {
	const memElementsSize = 32
	const bufferedMemCount = 10
	var owners []uint32
	var contents [bufferedMemCount][]uint32

	var start int = 0
	var end int = 4
	for end  <= len(ownersBytes) {
		owners = append(owners, BytesToI32(ownersBytes[start:end]))
		start = end
		end += 4
	}

	
	start = 0
	end = 4
	i := 0
	for end <= len(memoriesBytes) {
		if start % (memElementsSize * 4) == 0 && start != 0{
			i += 1
		}		
		contents[i] = append(contents[i], BytesToI32(memoriesBytes[start:end]))
		start = end
		end += 4
	}

	var memories []Memory
	for i := 0; i < len(owners); i++ {
		memory := Memory {
			Owner: owners[i],
			Content: contents[i],
		}
		memories = append(memories, memory)
	}

	return memories
}

func main() {
	// print(C.initialize(false))
	
	address := "0.0.0.0"
	port := 1337
	log.Printf("Listening on %s:%d\n", address, port)
	buffer := make([]byte, 1024)
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
		memoriesCount, owners, contents := ParseBuffer(buffer, size)
		fmt.Println("Number of Memories =", memoriesCount, "from", remoteaddr)
		memories := ParseMemories(memoriesCount, owners, contents)
		fmt.Println("Memories =", memories)
	}
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
