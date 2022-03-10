package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"io"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"
)

const (
	REDIS_ENDPOINT    = "http://127.0.0.1:6379"
	MAX_ALLOWED_PEERS = 5
)

type Distance struct {
	Peer     string
	Distance uint64
}

func GetCurrentRandomness() ([]byte, error) {
	var urls = []string{
		"https://api.drand.sh",
		"https://drand.cloudflare.com",
	}

	var chainHash, _ = hex.DecodeString("8990e7a9aaed2ffed73dbd7092123d6f289930540d7651336225dc172e51b2ce")
	c, err := client.New(
		client.From(http.ForURLs(urls, chainHash)...),
		client.WithChainHash(chainHash),
	)
	if err != nil {
		log.Fatal(err)
		return []byte{}, err
	}
	r, err := c.Get(context.Background(), 0)
	if err != nil {
		log.Fatal(err)
		return []byte{}, err
	}
	return r.Randomness(), nil
}

func GetDistance(a string, b string, salt string) uint64 {
	aHash := HashSHA256(a)
	bBytes := []byte(b)
	saltBytes := []byte(salt)
	sumBytes := make([]byte, 0)
	for i := 0; i < len(bBytes); i++ {
		sumBytes = append(sumBytes, bBytes[i]+saltBytes[i])
	}

	sumBytesHash := HashSHA256(string(sumBytes))
	distance := make([]byte, 0)
	for i := 0; i < len(aHash); i++ {
		distance = append(distance, aHash[i]^sumBytesHash[i])
	}

	return binary.BigEndian.Uint64(distance)
}

func GetPrivateSalt() [32]byte {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return HashSHA256(string(random.Uint64()))
}

func SendPeeringRequest(address string, port int) {
	buffer := make([]byte, 512)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	publickey, err := GetKey("pubkey")
	if err != nil {
		log.Fatal(err)
	}

	request := peering.Request{
		Publickey: publickey,
	}

	data, err := proto.Marshal(&request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[*] CLIENT: Sending messsage to %s:%d", address, port)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Fatal(err)
		return
	}

	conn.Close()
	return
}

func SendPeeringResponse(response peering.Response, address string, port int) {
	buffer := make([]byte, 512)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := proto.Marshal(&response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[*] CLIENT: Sending response to %s:%d with result %b", address, port, response.Result)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Fatal(err)
		return
	}

	conn.Close()
	return
}

// func GetCurrentPeers() []Distance {
// 		panic("todo :)")
// }

func SavePeerDistance(distance Distance) {
	pool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_ENDPOINT)
		},

		MaxIdle:     1024,
		IdleTimeout: 5 * time.Minute,
	}

	conn := pool.Get()
	defer conn.Close()

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(distance)
	conn.Do("SADD", "distances", buf.Bytes())
}

func RemovePeerDistance(distance Distance) {
	pool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_ENDPOINT)
		},

		MaxIdle:     1024,
		IdleTimeout: 5 * time.Minute,
	}

	conn := pool.Get()
	defer conn.Close()

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(distance)
	conn.Do("SREM", "distances", buf.Bytes())
}

func GetPeerDistances() []Distance {
	pool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_ENDPOINT)
		},

		MaxIdle:     1024,
		IdleTimeout: 5 * time.Minute,
	}

	conn := pool.Get()
	defer conn.Close()

	bs, _ := redis.Bytes(conn.Do("SMEMBERS", "distances"))
	bytesReader := bytes.NewReader(bs)

	var distances []Distance
	gob.NewDecoder(bytesReader).Decode(&distances)

	return distances
}

func EvaluatePeeringRequest(request peering.Request) peering.Response {
	var response peering.Response
	distances := GetPeerDistances()
	pubkey, _ := GetKey("pubkey")
	if len(distances) < MAX_ALLOWED_PEERS {
		proof := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		signature, checksum := Sign(proof)
		response = peering.Response{
			Result:    true,
			Proof:     proof,
			Signature: signature,
			Publickey: pubkey,
			Checksum: checksum,
		}
		myPubkey := fmt.Sprintf("%x", HashSHA256(response.Publickey))
		peerPubKey := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		privateSalt := fmt.Sprintf("%x", GetPrivateSalt())
		peerDistance := GetDistance(myPubkey, peerPubKey, privateSalt)
		distance := Distance{
			Peer:     request.Publickey,
			Distance: peerDistance,
		}
		SavePeerDistance(distance)
	} else {
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].Distance < distances[j].Distance
		})
		myPubkey := fmt.Sprintf("%x", HashSHA256(pubkey))
		peerPubKey := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		privateSalt := fmt.Sprintf("%x", GetPrivateSalt())
		peerDistance := GetDistance(myPubkey, peerPubKey, privateSalt)
		if distances[len(distances)-1].Distance > peerDistance {
			proof := fmt.Sprintf("%x", HashSHA256(request.Publickey))
			signature, checksum := Sign(proof)
			response = peering.Response{
				Result:    true,
				Proof:     proof,
				Signature: signature,
				Publickey: pubkey,
				Checksum: checksum,
			}
			distance := Distance{
				Peer:     request.Publickey,
				Distance: peerDistance,
			}
			RemovePeerDistance(distances[len(distances)-1])
			SavePeerDistance(distance)
		} else {
			response = peering.Response{
				Result:    false,
				Proof:     "",
				Signature: "",
				Publickey: pubkey,
				Checksum: "",
			}
		}
	}
	return response
}

func GenerateKeyPair() {
	priv, pub := GenerateRSAKeyPair()
	priv_pem := ExportRSAPrivateKey(priv)
	pub_pem := ExportRSAPublicKey(pub)
	fmt.Println(priv_pem, pub_pem)
}

func sendResponse(request peering.Request, conn *net.UDPConn, addr *net.UDPAddr) {
	response := EvaluatePeeringRequest(request)
	responseProto, _ := proto.Marshal(&response)
	_, err := conn.WriteToUDP(responseProto, addr)
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(conn *net.UDPConn) {
	buffer := make([]byte, 2048)
	for {
		size, remoteaddr, err := conn.ReadFromUDP(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
			break
		}

		log.Printf("Reading a request of size %d bytes from %v\n", size, remoteaddr)

		request := peering.Request{}
		err = proto.Unmarshal(buffer[:size], &request)
		if err != nil {
			log.Println(err.Error())
		}

		go sendResponse(request, conn, remoteaddr)
	}
}

func ListenForRequests(address string, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ser.Close()

	log.Printf("Listening on %s:%d\n", address, port)
	for {
		go handleRequest(ser)
	}
}

func main() {
	fmt.Println("Hello")

	GenerateKeyPair()

	randomness, _ := GetCurrentRandomness()
	fmt.Printf("Current Randomness: %x\n", randomness)

	privateSalt := GetPrivateSalt()
	fmt.Printf("Private Salt: %x\n", privateSalt)
	a := "877133ac2143ac542a2f0e7c415705770b9d47dc8f13d0b7f2c7346ae52eee24"
	b := "362c3452adbf9616c33a42dd8ae5cf651c197cc663807b1ca7d7a2006229fa29"
	salt := "a2054fccdc58815afb604f5742df7def70a279f3e86f0142365d9796cf14091f"
	fmt.Println("Distance:", GetDistance(a, b, salt))

	request := peering.Request{
		Publickey: "362c3452adbf9616c33a42dd8ae5cf651c197cc663807b1ca7d7a2006229fa29",
	}
	response := EvaluatePeeringRequest(request)
	fmt.Println(response)
}
