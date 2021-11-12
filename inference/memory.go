package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
	"os"
)

type Memory struct {
	Owner	uint32
	Content []uint32
}

func I32ToBytes(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func BytesToI32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func AppendToBytes(buffer *[]byte, value []byte) {
	for i := 0; i < 4; i++ {
		*buffer = append(*buffer, value[i])
	}
}

func PrepareMessage(memories []Memory) []byte {
	const memoryElementsCount = 32
	const bufferedMemCount = 10
	var memoriesCount int = len(memories)

	buffer := make([]byte, 0)
	
	AppendToBytes(&buffer, I32ToBytes(uint32(memoriesCount)))
	
	var contents [bufferedMemCount][memoryElementsCount]uint32
	for i := 0; i < memoriesCount; i++ {
		AppendToBytes(&buffer, I32ToBytes(memories[i].Owner))
		for j := 0; j < memoryElementsCount; j++ {
			contents[i][j] = memories[i].Content[j]
		}
	}

	for i := 0; i < memoriesCount; i++ {
		for j := 0; j < memoryElementsCount; j++ {
			AppendToBytes(&buffer, I32ToBytes(contents[i][j]))
		}
	}

	return buffer
}

func SaveMemories(memories []Memory) {
	var buffer []byte = PrepareMessage(memories)
	now := time.Now()
	timestamp := now.Unix()	
	filename := fmt.Sprintf("./memories/%d.mem", timestamp)
	f, err := os.Create(filename)
	
	if err != nil {
		log.Fatal(err)
	}
	
	f.Write(buffer)
	f.Close()	
}

func LoadMemories(timestamp int) []Memory {
	filename := fmt.Sprintf("./memories/%d.mem", timestamp)
	
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	memoriesCount, owners, contents := ParseBuffer(buffer)
	memories := ParseMemories(memoriesCount, owners, contents)
	
	return memories
}

func ParseBuffer(buffer []byte) (uint32, []byte, []byte) {
	const memoryElementsCount = 32
	var memoriesCount uint32
	memoriesCount = BytesToI32(buffer[0:4])
	var ownersBytes []byte = buffer[4:4+memoriesCount*4]
	var contentsBytes []byte = buffer[4+memoriesCount*4:4+memoriesCount*4+memoriesCount*memoryElementsCount*4]
	return memoriesCount, ownersBytes, contentsBytes
}

func ParseMemories(memoriesCount uint32, ownersBytes []byte, memoriesBytes []byte) []Memory {
	const memoryElementsCount = 32
	const bufferedMemCount = 10
	var owners []uint32
	var contents [bufferedMemCount][]uint32

	var start int = 0
	var end int = 4
	for end  <= len(ownersBytes) {
		owners = append(owners, BytesToI32(ownersBytes[start:end]))
		start = end
		end += 4
	}

	start = 0
	end = 4
	i := 0
	for end <= len(memoriesBytes) {
		if start % (memoryElementsCount * 4) == 0 && start != 0{
			i += 1
		}		
		contents[i] = append(contents[i], BytesToI32(memoriesBytes[start:end]))
		start = end
		end += 4
	}

	var memories []Memory
	for i := 0; i < len(owners); i++ {
		memory := Memory {
			Owner: owners[i],
			Content: contents[i],
		}
		memories = append(memories, memory)
	}

	return memories
}

func SendMemories() {
	// TODO(ahmed): change address to C.get_random_peer()
	address := "0.0.0.0"
	port := 1337

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Println(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}

	memories := GenerateMemories()
	memoriesBytes := string(PrepareMessage(memories))
	
	log.Printf("Sending memories as a batch to %s:%d", address, port)
	fmt.Fprintf(conn, memoriesBytes)

	conn.Close()
}

func GenerateMemories() []Memory {
	memories := []Memory{
		{
			Owner: 7,
			Content: []uint32{
				1,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				9,
			},
		},
		{
			Owner: 2237794909,
			Content: []uint32{
				1350226788,
				2035283319,
				3328873732,
				2013792798,
				2237794909,
				2888525616,
				3046543999,
				3617579486,
				1971669918,
				2513307423,
				1245590728,
				2987079650,
				1876808204,
				1925279845,
				1956089770,
				1705405968,
				1038792516,
				3256599353,
				1502246338,
				3616156870,
				4053352369,
				1326285247,
				1456192919,
				1838165923,
				1918116757,
				1044999207,
				3595028482,
				3453924757,
				3064481954,
				1470385246,
				3946967500,
				1719886239,
			},
		},
	}
	return memories
}
