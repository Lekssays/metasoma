package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"sort"
	"time"
)

const (
	REDIS_ENDPOINT    = "http://127.0.0.1:6379"
	MAX_ALLOWED_PEERS = 5
)

type Distance struct {
	Publickey string
	Address   string
	Port      uint32
	Value     uint64
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
			Checksum:  checksum,
			Type:      "PEERING",
		}
		myPubkey := fmt.Sprintf("%x", HashSHA256(response.Publickey))
		peerPubKey := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		privateSalt := fmt.Sprintf("%x", GetPrivateSalt())
		peerDistance := GetDistance(myPubkey, peerPubKey, privateSalt)
		distance := Distance{
			Publickey: request.Publickey,
			Address:   request.Address,
			Port:      request.Port,
			Value:     peerDistance,
		}
		SavePeerDistance(distance)
	} else {
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].Value < distances[j].Value
		})
		myPubkey := fmt.Sprintf("%x", HashSHA256(pubkey))
		peerPubKey := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		privateSalt := fmt.Sprintf("%x", GetPrivateSalt())
		peerDistance := GetDistance(myPubkey, peerPubKey, privateSalt)
		if distances[len(distances)-1].Value > peerDistance {
			proof := fmt.Sprintf("%x", HashSHA256(request.Publickey))
			signature, checksum := Sign(proof)
			response = peering.Response{
				Result:    true,
				Proof:     proof,
				Signature: signature,
				Publickey: pubkey,
				Checksum:  checksum,
				Type:      "PEERING",
			}
			distance := Distance{
				Publickey: request.Publickey,
				Address:   request.Address,
				Port:      request.Port,
				Value:     peerDistance,
			}
			RemovePeerDistance(distances[len(distances)-1])
			SavePeerDistance(distance)
		} else {
			response = peering.Response{
				Result:    false,
				Proof:     "",
				Signature: "",
				Publickey: pubkey,
				Checksum:  "",
				Type:      "PEERING",
			}
		}
	}
	return response
}
