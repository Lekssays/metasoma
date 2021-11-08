package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
	"log"
)

type Memory struct {
	Owner uint32
	Hosts []uint32
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

func GenerateMemories() []Memory {
	memories := []Memory{
		{
			Owner: 3290510358,
			Hosts: []uint32{
				3103580358,
				3210358905,
				3210358905,
				1033290558,
				2989952514,
				3210358905,
				1033290558,
				3103580358,
				2652254021
			},
		},
		{
			Owner: 2652254021,
			Hosts: []uint32{
				2989952514,
				3210358905,
				1033290558
			},
		},
		{
			Owner: 2989952514,
			Hosts: []uint32{
				3103580358,
				2652254021,
				3290510358,
				2989952514,
				3210358905,
				1033290558
			},
		},
	}
	return memories
}
