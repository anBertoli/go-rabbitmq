package main

import (
	"flag"
	"log"
)

const help = "choose between 'consumer' or 'producer'"

func main() {
	mode := flag.String("mode", "producer", help)
	flag.Parse()

	switch *mode {
	case "producer":
		producer()
	case "consumer":
		consumer()
	default:
		log.Fatalf(help)
	}
}
