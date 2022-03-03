package main

import (
	"context"
	"fmt"

	"encoding/hex"

	"github.com/drand/drand/client"
	"github.com/drand/drand/client/http"
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

func GetDistance(a string, b string) int {
	panic("todo :)")
}

func GetPotentialPeers() {
	panic("todo :)")
}

func getPrivateSalt() (string, error) {
	key := "privkey"
	privkey, err := GetKey(key)
	if err != nil {
		return "", err
	}
	return HashSHA256(privkey), nil
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
	fmt.Println(randomness)

	privateSalt, _ := getPrivateSalt()
	fmt.Println(privateSalt)
}
