package main

import (
	"flag"
	"log"
)

const help = "choose between 'publisher' or 'subscriber'"

func main() {
	mode := flag.String("mode", "", help)
	flag.Parse()

	switch *mode {
	case "publisher":
		publisher()
	case "subscriber":
		subscriber()
	default:
		log.Fatalf(help)
	}
}
