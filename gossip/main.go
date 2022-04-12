package main

import (
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	args := os.Args[1:]

	if args[0] == "server" {
		for {
			timer := time.After(2 * time.Second)

			wg.Add(1)
			go RunServer(&wg)

			wg.Wait()

			<-timer
		}
	} else if args[0] == "client" {
		RunClient()
	}
}
