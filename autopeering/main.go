package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Lekssays/metasoma/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("Starting Autopeering Service :)...")

	if _, err := os.Stat("./" + DISCOVERY_ADDRESS + "_pubkey.pem"); errors.Is(err, os.ErrNotExist) {
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
				if rtype.Response.Purpose == peering.Purpose_PONG {
					log.Printf("Receiving %s response from %v", rtype.Response.Purpose, remoteaddr)
				} else if rtype.Response.Purpose == peering.Purpose_PEERING {
					log.Printf("Receiving %s response with result %s from %v", rtype.Response.Purpose, strconv.FormatBool(rtype.Response.Result), remoteaddr)
					EvaluateResponse(rtype.Response)
				} else if rtype.Response.Purpose == peering.Purpose_GOSSIP {
					log.Printf("Receiving %s response from %v", rtype.Response.Purpose, remoteaddr)
					for i := 0; i < len(rtype.Response.Peers); i++ {
						peer := rtype.Response.Peers[i]
						log.Printf("Sending %s request to %s:%d", rtype.Response.Purpose, peer.Address, int(peer.Port))
						SendRequest(peering.Purpose_PEERING, peer.Address, int(peer.Port))
					}
				}
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
		var wg sync.WaitGroup
		sent := make(map[string]bool)
		for {
			timer := time.After(1 * time.Second)

			// wg.Add(1)
			// go CheckLiveness(&wg)

			wg.Add(1)
			go GossipPeers(&wg, sent)

			wg.Wait()

			<-timer
		}
	} else if args[0] == "simulator" {
		now := time.Now().Unix()
		SIMULATION_PERIOD := 100
		for time.Now().Unix()-now <= int64(SIMULATION_PERIOD) {
			peers, _ := GetCurrentPeers()
			neighbors, _ := GetPeersDistances()
			entry := fmt.Sprintf("peers: %v", peers)
			WriteLog(entry)
			entry = fmt.Sprintf("neighbors: %v", neighbors)
			WriteLog(entry)
			time.Sleep(1 * time.Second)
		}
	}
}
