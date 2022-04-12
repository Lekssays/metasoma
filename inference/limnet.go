package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println(C.initialize(false))

	// Example of handling a packet
	filePath := "./data/192.168.1.252.csv"
	records := ReadCSV(filePath)
	packets := ParsePackets(records)
	C.on_packet_received(C.uint(packets[0].srcIp), C.uint(packets[0].dstIp), (*C.float)(&packets[0].features[0]))

	m := sync.Mutex{}
	m.Lock()
	m.Lock()
}
