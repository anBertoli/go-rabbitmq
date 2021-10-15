package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

const (
	helpMode = "choose between 'client' or 'server'"
	rpcQueue = "rpc-queue"
)

func main() {
	rand.Seed(time.Now().Unix())

	mode := flag.String("mode", "", helpMode)
	flag.Parse()

	switch *mode {
	case "client":
		x()
	case "server":

	default:
		log.Fatalf(helpMode)
	}
}
