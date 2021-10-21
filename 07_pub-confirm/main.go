package main

import (
	"math/rand"
	"time"
)

const (
	confirmationQueue = "pub-confirm"
)

func main() {
	rand.Seed(time.Now().Unix())
	start()
}
