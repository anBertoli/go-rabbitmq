package main

import (
	"flag"
	"log"
)

const help = "choose between 'worker' or 'producer'"

func main() {
	mode := flag.String("mode", "producer", help)
	task := flag.String("task", "", help)
	flag.Parse()

	switch *mode {
	case "producer":
		parsedTask := taskFromArgs(*task)
		producer(parsedTask)
	case "worker":
		worker()
	default:
		log.Fatalf(help)
	}
}
