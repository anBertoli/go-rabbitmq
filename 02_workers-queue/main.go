package main

import (
	"flag"
	"log"
)

const help = "choose between 'worker' or 'producer'"

func main() {
	mode := flag.String("mode", "producer", help)
	flag.Parse()

	switch *mode {
	case "producer":
		producer()
	case "worker":
		worker()
	default:
		log.Fatalf(help)
	}
}
