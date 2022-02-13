package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var cache = make(map[int]Memory)

func Set(memory Memory) {
	cache[memory.From] = memory
}

func Get(owner uint32) Memory {
	return cache[owner]
}

func Store(s []byte, owner uint32) {
	filename := fmt.Sprintf("./memories/%d.mem", owner)
	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	f.Write(s)
	f.Close()
}

func Read(owner uint32) []byte {
	filename := fmt.Sprintf("./memories/%d.mem", owner)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func Backup() {
	for _, memory := range cache {
		data := Encode(memory)
		Store(data, memory.Owner)
	}
}

func main() {
	memories := GenerateMemories()
	for _, m := range memories {
		Set(m)
	}

	ids := []int{3290510358, 2652254021, 2989952514}
	for _, id := range ids {
		memory := Get(id)
		fmt.Println(memory)
	}

	data := Encode(memories[1])
	Store(data, memories[1].From)

	readData := Read(memories[1].From)
	readMemory := Decode(readData)
	fmt.Println(readMemory)

	Backup()
}
