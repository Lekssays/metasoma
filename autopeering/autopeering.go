package main

import (
	"context"
	"fmt"

	"encoding/binary"
	"encoding/hex"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"

	"math/rand"
	"time"
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
		return []byte{}, err
	}
	r, err := c.Get(context.Background(), 0)
	if err != nil {
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
		sumBytes = append(sumBytes, bBytes[i] + saltBytes[i])
	}
	
	sumBytesHash := HashSHA256(string(sumBytes))
	distance := make([]byte, 0)
	for i := 0; i < len(aHash); i++ {
		distance = append(distance, aHash[i] ^ sumBytesHash[i])
	}
	
	return binary.BigEndian.Uint64(distance)
}

func GetPotentialPeers() {
	panic("todo :)")
}

func getPrivateSalt() [32]byte {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return HashSHA256(string(random.Uint64()))
}

func SendPeeringRequest() {
	panic("todo :)")
}

func AcceptPeeringRequest() {
	panic("todo :)")
}

func GenerateKeyPair() {
	priv, pub := GenerateRSAKeyPair()
	priv_pem := ExportRSAPrivateKey(priv)
	pub_pem := ExportRSAPublicKey(pub)
	fmt.Println(priv_pem, pub_pem)
}

func main() {
	fmt.Println("Hello")
	
	GenerateKeyPair()

	randomness, _ := GetCurrentRandomness()
	fmt.Printf("Current Randomness: %x\n", randomness)

	privateSalt := getPrivateSalt()
	fmt.Printf("Private Salt: %x\n", privateSalt)
	a := "877133ac2143ac542a2f0e7c415705770b9d47dc8f13d0b7f2c7346ae52eee24"
	b := "362c3452adbf9616c33a42dd8ae5cf651c197cc663807b1ca7d7a2006229fa29"
	salt := "a2054fccdc58815afb604f5742df7def70a279f3e86f0142365d9796cf14091f"
	fmt.Println("Distance:", GetDistance(a, b, salt))
}
