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

			payload := peering.Payload{}
			err = proto.Unmarshal(buffer[:size], &payload)
			if err != nil {
				log.Println(err.Error())
			}

			switch rtype := payload.Type.(type) {
			case *peering.Payload_Response:
				log.Printf("Receiving %s response with result %s from %v", rtype.Response.Purpose, strconv.FormatBool(rtype.Response.Result), remoteaddr)
			case *peering.Payload_Request:
				log.Printf("Receiving %s request from %v", rtype.Request.Purpose, remoteaddr)
				receivingAddress := net.UDPAddr{
					Port: int(rtype.Request.Port),
					IP:   net.ParseIP(remoteaddr.IP.String()),
				}
				go SendResponse(rtype.Request, ser, &receivingAddress)
			}
		}
	} else if args[0] == "client" {
		CheckLivness()
	}
}
