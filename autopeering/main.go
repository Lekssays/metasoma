package main

import (
	"fmt"
	"github.com/Lekssays/ADeLe/autopeering/protos/peering"
)

func main() {
	fmt.Println("Autopeering :)")

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
