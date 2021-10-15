package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

const (
	helpMode          = "choose between 'publisher' or 'subscriber'"
	helpBind          = "format '<facility>.<severity>'"
	logsTopicExchange = "logs-topic-exchange"
)

func main() {
	mode := flag.String("mode", "", helpMode)
	bind := flag.String("bind", "", helpBind)
	flag.Parse()

	rand.Seed(time.Now().Unix())

	switch *mode {
	case "publisher":
		publisher()
	case "subscriber":
		if !validateBinding(*bind) {
			log.Fatalf("invalid binding: '%s'", *bind)
		}
		subscriber(*bind)
	default:
		log.Fatalf(helpMode)
	}
}
