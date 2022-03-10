package main

import (
	"bufio"
	"fmt"
	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
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
	buffer := make([]byte, 512)

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

	log.Printf("Sending %s request to %s:%d", request.Type, address, port)
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
			Proof:     "",
			Signature: "",
			Publickey: pubkey,
			Checksum:  "",
			Type:      "PONG",
		}
	}

	responseProto, _ := proto.Marshal(&response)
	_, err := conn.WriteToUDP(responseProto, addr)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	args := os.Args[1:]

	if args[0] == "server" {
		addr := net.UDPAddr{
			Port: DISCOVERY_PORT,
			IP:   net.ParseIP(DISCOVERY_ADDRESS),
		}
		ser, err := net.ListenUDP("udp", &addr)
		if err != nil {
			log.Println(err)
			return
		}
		defer ser.Close()

		log.Printf("Listening on %s:%d\n", DISCOVERY_ADDRESS, DISCOVERY_PORT)
		buffer := make([]byte, 2048)
		for {
			size, remoteaddr, err := ser.ReadFromUDP(buffer)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
				break
			}

			request := peering.Request{}
			err = proto.Unmarshal(buffer[:size], &request)
			if err != nil {
				log.Println(err.Error())
			}

			log.Printf("Receiving %s request from %v\n", request.Type, remoteaddr)

			endpoint := strings.Split(remoteaddr.String(), ":")
			if endpoint[0] == request.Address {
				receivingAddress := net.UDPAddr{
					Port: int(request.Port),
					IP:   net.ParseIP(request.Address),
				}

				go SendResponse(request, ser, &receivingAddress)
			}
		}
	} else if args[0] == "client" {
		SendRequest("PING", DISCOVERY_ADDRESS, DISCOVERY_PORT)
	}
}
