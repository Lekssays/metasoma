package main

import (
	"fmt"
	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting Autopeering Service :)...")

	if _, err := os.Stat("pubkey.pem"); errors.Is(err, os.ErrNotExist) {
		GenerateKeyPair()
	}
	
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
