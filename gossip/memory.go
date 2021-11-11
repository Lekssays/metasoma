package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"log"
)

type Memory struct {
	Owner	uint32
	Content [32]uint32
}

func Encode(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("uncompressed size (bytes): ", len(buf.Bytes()))

	return buf.Bytes()
}

func Compress(s []byte) []byte {
	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	log.Println("compressed size (bytes): ", len(zipbuf.Bytes()))
	return zipbuf.Bytes()
}

func Decompress(s []byte) []byte {
	rdr, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := ioutil.ReadAll(rdr)

	if err != nil {
		log.Fatal(err)
	}

	rdr.Close()
	log.Println("uncompressed size (bytes): ", len(data))

	return data
}

func Decode(s []byte) Memory {
	memory := Memory{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&memory)

	if err != nil {
		log.Fatal(err)
	}

	return memory
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
	const compressedMemSize = 32
	const maxMemCount = 10
	var memoriesCount int = len(memories)

	buffer := make([]byte, 0)
	
	AppendToBytes(&buffer, I32ToBytes(uint32(memoriesCount)))
	
	var contents [maxMemCount][compressedMemSize]uint32
	for i := 0; i < memoriesCount; i++ {
		AppendToBytes(&buffer, I32ToBytes(memories[i].Owner))
		for j := 0; j < compressedMemSize; j++ {
			contents[i][j] = memories[i].Content[j]
		}
	}

	for i := 0; i < memoriesCount; i++ {
		for j := 0; j < compressedMemSize; j++ {
			AppendToBytes(&buffer, I32ToBytes(contents[i][j]))
		}
	}

	return buffer
}

func GenerateMemories() []Memory {
	memories := []Memory{
		{
			Owner: 7,
			Content: [32]uint32{
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
			Content: [32]uint32{
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
