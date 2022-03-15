package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("Starting Autopeering Service :)...")

	if _, err := os.Stat("./pubkey.pem"); errors.Is(err, os.ErrNotExist) {
		GenerateKeyPair()
		LoadBasePeers()
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

			// todo(ahmed): there is an elegant way of doing it with oneof in proto definition
			if request.Type != "PEERING" || request.Type != "PING" {
				response := peering.Response{}
				err = proto.Unmarshal(buffer[:size], &response)
				if err != nil {
					log.Println(err.Error())
				}
				fmt.Println(response)
				log.Printf("Receiving %s response with result %s from %v", response.Type, strconv.FormatBool(response.Result), remoteaddr)
			} else {
				log.Printf("Receiving %s request from %v", request.Type, remoteaddr)
				receivingAddress := net.UDPAddr{
					Port: int(request.Port),
					IP:   net.ParseIP(remoteaddr.IP.String()),
				}
				go SendResponse(request, ser, &receivingAddress)
			}

		}
	} else if args[0] == "client" {
		CheckLivness()
	}
}
