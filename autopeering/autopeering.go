package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

const (
	MAX_ALLOWED_PEERS = 5
)

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

func SavePeerDistance(distance peering.Distance) (bool, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_SERVER,
		Password: "",
		DB:       0,
	})

	distanceString := proto.MarshalTextString(&distance)

	err := rdb.SAdd(ctx, "distances", distanceString).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemovePeerDistance(distance peering.Distance) (bool, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_SERVER,
		Password: "",
		DB:       0,
	})

	distanceString := proto.MarshalTextString(&distance)
	err := rdb.SRem(ctx, "distances", distanceString).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetPeerDistances() ([]peering.Distance, error) {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_SERVER,
		Password: "",
		DB:       0,
	})

	distancesStrings, err := rdb.SMembers(ctx, "distances").Result()
	if err != nil {
		return []peering.Distance{}, err
	}

	var distances []peering.Distance
	for i := 0; i < len(distancesStrings); i++ {
		var distance peering.Distance
		err = proto.UnmarshalText(distancesStrings[i], &distance)
		if err != nil {
			log.Println(err.Error())
		}
		distances = append(distances, distance)
	}

	return distances, nil
}

func EvaluatePeeringRequest(request *peering.Request) peering.Response {
	var response peering.Response
	distances, _ := GetPeerDistances()
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
			Purpose:   peering.Purpose_PEERING,
		}
		myPubkey := fmt.Sprintf("%x", HashSHA256(response.Publickey))
		peerPubKey := fmt.Sprintf("%x", HashSHA256(request.Publickey))
		privateSalt := fmt.Sprintf("%x", GetPrivateSalt())
		peerDistance := GetDistance(myPubkey, peerPubKey, privateSalt)
		distance := peering.Distance{
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
				Purpose:   peering.Purpose_PEERING,
			}
			distance := peering.Distance{
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
				Proof:     "null",
				Signature: "null",
				Publickey: pubkey,
				Checksum:  "null",
				Purpose:   peering.Purpose_PEERING,
			}
		}
	}
	return response
}

func GetCurrentPeers() []string {
	endpoints := make([]string, 0)
	distances, _ := GetPeerDistances()
	for i := 0; i < len(distances); i++ {
		endpoints = append(endpoints, fmt.Sprintf("%s:%d", distances[i].Address, distances[i].Port))
	}
	return endpoints
}

func LoadBasePeers() {
	// todo(ahmed): Load peers from a generated binary (nodes managed by LiMNet admins)
	var distance peering.Distance
	if DISCOVERY_ADDRESS == "peer0.limnet.io" {
		distance = peering.Distance{
			Publickey: "",
			Address:   "peer1.limnet.io",
			Port:      1337,
			Value:     100,
		}
	} else if DISCOVERY_ADDRESS == "peer1.limnet.io" {
		distance = peering.Distance{
			Publickey: "",
			Address:   "peer0.limnet.io",
			Port:      1337,
			Value:     100,
		}
	}
	SavePeerDistance(distance)
}
