package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
)

var (
	DISCOVERY_ADDRESS = os.Getenv("DISCOVERY_ADDRESS")
	DISCOVERY_PORT    = getPort()
)

func getPort() int {
	val := os.Getenv("DISCOVERY_PORT")
	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

func SendRequest(rtype string, address string, port int) {
	buffer := make([]byte, 2048)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	publickey, err := GetKey("pubkey")
	if err != nil {
		log.Fatal(err)
	}

	request := peering.Request{
		Publickey: publickey,
		Address:   DISCOVERY_ADDRESS,
		Port:      uint32(DISCOVERY_PORT),
		Type:      rtype,
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sending %s request from %s:%d to %s:%d", request.Type, request.Address, request.Port, address, port)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Fatal(err)
	}

	conn.Close()
}

func SendResponse(request peering.Request, conn *net.UDPConn, addr *net.UDPAddr) {
	var response peering.Response
	pubkey, _ := GetKey("pubkey")
	if request.Type == "PEERING" {
		response = EvaluatePeeringRequest(request)
	} else if request.Type == "PING" {
		response = peering.Response{
			Result:    false,
			Proof:     "null",
			Signature: "null",
			Publickey: pubkey,
			Checksum:  "null",
			Type:      "PONG",
		}
	}

	responseProto, _ := proto.Marshal(&response)
	_, err := conn.WriteToUDP(responseProto, addr)
	if err != nil {
		log.Println(err)
	}
}

func CheckLivness() {
	for {
		endpoints := GetCurrentPeers()
		if len(endpoints) > 0 {
			for i := 0; i < len(endpoints); i++ {
				endpoint := strings.Split(endpoints[i], ":")
				address := endpoint[0]
				port, err := strconv.Atoi(endpoint[1])
				if err != nil {
					log.Fatal(err)
				}
				if address != DISCOVERY_ADDRESS {
					SendRequest("PING", address, port)
				}
			}
		}
		time.Sleep(60 * time.Second)
	}
}
