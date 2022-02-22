package main

import (
	"github.com/Lekssays/ADeLe/gossip/proto/message"
)

func GetMergedMemory(memories []message.Memory) message.Memory {
	// todo(ahmed): implement getting merged memory
	return message.Memory{
		From:   1350226788,
		Target: 3328873732,
		Checksum:  "997b28e69f45f7bcdaca04dfe9562b3ef59f3dd3064220f0410a9c03f0ddcc82",
		Signature: "OTk3YjI4ZTY5ZjQ1ZjdiY2RhY2EwNGRmZTk1NjJiM2VmNTlmM2RkMzA2NDIyMGYwNDEwYTljMDNmMGRkY2M4Mg==",
		Content: []uint32{1470385246, 1038792516, 1956089770, 1876808204, 1350226788},
	}
}

func PrepareMessage(memories []message.Memory, peeringProofs []string) message.Message {
	mergedMemory := GetMergedMemory(memories)
	parents := []*message.Memory{&memories[0], &memories[1]}
	return message.Message{
		MergedMemory: &mergedMemory,
		Parents: parents,
		PeeringProofs: peeringProofs,
	}
}

func GenerateDummyMemories() []message.Memory {
	memories := []message.Memory{
		{
			From:   1350226788,
			Target: 3328873732,
			Checksum:  "997b28e69f45f7bcdaca04dfe9562b3ef59f3dd3064220f0410a9c03f0ddcc82",
			Signature: "OTk3YjI4ZTY5ZjQ1ZjdiY2RhY2EwNGRmZTk1NjJiM2VmNTlmM2RkMzA2NDIyMGYwNDEwYTljMDNmMGRkY2M4Mg==",
			Content: []uint32{1470385246, 1038792516, 1956089770, 1876808204, 1350226788},
		},
		{
			From: 1705405968,
			Target: 1159684051,
			Checksum: "ac9b54934b516627f9e44d2965c6e4f1bb3f8aeafff24ad7178c2ee3f07cbc21",
			Signature: "YWM5YjU0OTM0YjUxNjYyN2Y5ZTQ0ZDI5NjVjNmU0ZjFiYjNmOGFlYWZmZjI0YWQ3MTc4YzJlZTNmMDdjYmMyMQ==",
			Content: []uint32{1326285247, 1118116757, 2295028482},
		},
		{
			From:   2237794909,
			Target: 1971669918,
			Content: []uint32{1350226788, 2035283319, 3328873732, 2013792798, 2237794909, 3064481954},
			Checksum:  "35cdabe030a7eabaa6d8507b42b09b74077eacb10881af5d0cdcd671fda51829",
			Signature: "MzVjZGFiZTAzMGE3ZWFiYWE2ZDg1MDdiNDJiMDliNzQwNzdlYWNiMTA4ODFhZjVkMGNkY2Q2NzFmZGE1MTgyOQ==",
		},
	}
	return memories
}

func GenerateDummyPeeringProofs() []string {
	return []string{
		"MzVjZGFiZTAzMGE3ZWFiYWE2ZDg1MDdiNDJiMDliNzQwNzdlYWNiMTA4ODFhZjVkMGNkY2Q2NzFmZGE1MTgyOQ==",
		"TXpWalpHRmlaVEF6TUdFM1pXRmlZV0UyWkRnMU1EZGlOREppTURsaU56UXdOemRsWVdOaU1UQTRPREZoWmpWa01HTmtZMlEyTnpGbVpHRTFNVGd5T1E9PQ==",
		"V2FscEhSbWxhVkVGNlRVZEZNMXBYUm1sWlYwVXlXa1JuTVUxRVpHbE9SRXBwVFVSc2E=",
	}
}