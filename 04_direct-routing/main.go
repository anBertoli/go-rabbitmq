package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

const (
	helpMode            = "choose between 'publisher' or 'subscriber'"
	helpSevs            = "choose one or more of 'info', 'warn' or 'error' (format: info-warn-error)"
	logsRoutingExchange = "logs-routing"
)

func main() {
	mode := flag.String("mode", "", helpMode)
	sev := flag.String("sevs", "", helpSevs)
	flag.Parse()

	rand.Seed(time.Now().Unix())

	switch *mode {
	case "publisher":
		publisher()
	case "subscriber":
		sevs, ok := validateSeverities(*sev)
		if !ok {
			log.Fatalf("invalid severity: '%s'", *sev)
		}
		subscriber(sevs)
	default:
		log.Fatalf(helpMode)
	}
}
