package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -llimnet -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

func main() {
	print(C.initialize(false))
	
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
		num_entries = size / (4 + C.compressed_memory_size())
		on_memories_received(&buffer[0], &buffer[4*num_entries], num_entries)
	}
}

func send_memories(conn *net.UDPConn) {
	peer := net.UDPAddr{
		Port: 1337,
		IP:  C.get_random_peer(),
	}
	num_entries = 10
	buffer := make([]byte, num_entries*(4+C.compressed_memory_size()))
	num_entries = C.get_memories_to_share(&buffer[0], &buffer[4*num_entries], num_entries)
	_, err := conn.WriteToUDP(buffer.slice(0, num_entries*(4+C.compressed_memory_size())), peer)
	if err != nil {
		log.Println(err)
	}
}
