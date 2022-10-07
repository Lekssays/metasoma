package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/Lekssays/metasoma/autopeering/protos/peering"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
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

func SendRequest(purpose peering.Purpose, address string, port int) {
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

	requestUUID := uuid.New().String()
	proof := fmt.Sprintf("%x", HashSHA256(GenerateProof()))
	payload := peering.Payload{
		Type: &peering.Payload_Request{
			Request: &peering.Request{
				Publickey: publickey,
				Address:   DISCOVERY_ADDRESS,
				Port:      uint32(DISCOVERY_PORT),
				Purpose:   purpose,
				Uuid:      requestUUID,
				Proof:     proof,
			},
		},
	}

	data, err := proto.Marshal(&payload)
	if err != nil {
		log.Fatal(err)
	}

	request := peering.Request{
		Publickey: publickey,
		Address:   DISCOVERY_ADDRESS,
		Port:      uint32(DISCOVERY_PORT),
		Purpose:   purpose,
		Uuid:      requestUUID,
		Proof:     proof,
	}

	_, err = SaveRequest(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sending %s request to %s:%d\n", purpose.String(), address, port)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Fatal(err)
	}

	conn.Close()
}

func SendResponse(request *peering.Request, conn *net.UDPConn, addr *net.UDPAddr) {
	var response peering.Response
	pubkey, _ := GetKey("pubkey")
	if request.Purpose == peering.Purpose_PEERING {
		response = EvaluatePeeringRequest(request)
	} else if request.Purpose == peering.Purpose_PING || request.Purpose == peering.Purpose_GOSSIP {
		peers, _ := GetCurrentPeers()
		response = peering.Response{
			Result:    false,
			Proof:     "null",
			Signature: "null",
			Publickey: pubkey,
			Checksum:  "null",
			Purpose:   peering.Purpose_GOSSIP,
			Uuid:      request.Uuid,
			Peers:     peers,
		}
	}

	payload := peering.Payload{
		Type: &peering.Payload_Response{
			Response: &response,
		},
	}

	payloadProto, _ := proto.Marshal(&payload)
	log.Printf("Sending %s response to %s:%d\n", response.Purpose, addr.IP.String(), addr.Port)
	_, err := conn.WriteToUDP(payloadProto, addr)
	if err != nil {
		log.Println(err)
	}
}

// func CheckLiveness(wg *sync.WaitGroup) {
// 	endpoints := GetCurrentPeers()
// 	if len(endpoints) > 0 {
// 		for i := 0; i < len(endpoints); i++ {
// 			endpoint := strings.Split(endpoints[i], ":")
// 			address := endpoint[0]
// 			port, err := strconv.Atoi(endpoint[1])
// 			if err != nil {
// 				log.Fatal(err)
// 				wg.Done()
// 			}
// 			if address != DISCOVERY_ADDRESS {
// 				SendRequest(peering.
// 			}
// 		}
// 	}
// 	time.Sleep(60 * time.Second)
// 	wg.Done()
// }

// func GossipPeers(wg *sync.WaitGroup, sent map[string]bool) {
// 	peers := make([]peering.Peer, 0)
// 	distances, _ := GetPeersDistances()
// 	for i := 0; i < len(distances); i++ {
// 		peer := peering.Peer{
// 			Publickey: distances[i].Publickey,
// 			Address:   distances[i].Address,
// 			Port:      distances[i].Port,
// 		}
// 		if !sent[peer.Publickey] {
// 			sent[peer.Publickey] = true
// 			peers = append(peers, peer)
// 		}
// 	}
// 	time.Sleep(1 * time.Second)
// 	wg.Done()
// }

func WriteLog(content string) {
	address := os.Getenv("DISCOVERY_ADDRESS")
	f, err := os.OpenFile(address+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(content)
}

func GossipPeers(wg *sync.WaitGroup, sent map[string]bool) {
	peers, _ := GetCurrentPeers()
	for i := 0; i < len(peers); i++ {
		if peers[i].Address != DISCOVERY_ADDRESS && !sent[peers[i].Address] {
			SendRequest(peering.Purpose_GOSSIP, peers[i].Address, int(peers[i].Port))
		}
	}
	// time.Sleep(2 * time.Second)
	wg.Done()
}
