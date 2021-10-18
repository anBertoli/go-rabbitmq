package main

import (
	"math/rand"
	"time"
)

const (
	confirmationQueue = "confirmation-queue"
)

func main() {
	rand.Seed(time.Now().Unix())
	start()
}
