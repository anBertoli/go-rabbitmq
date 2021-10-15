package main

import (
	"math/rand"
	"strings"
)

// List of possible severity values.
const (
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

// List of all possible sources.
const (
	NGINX = "nginx"
	CRON  = "cron"
	SSHD  = "sshd"
)

// Wildcards.
const (
	STAR = "*"
	HASH = "#"
)

// Validate subscriber bindings.
func validateBinding(bind string) bool {
	parts := strings.Split(bind, ".")
	if len(parts) != 2 {
		return false
	}

	ok := false
	for _, src := range []string{NGINX, CRON, SSHD, STAR, HASH} {
		if parts[0] == src {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}

	ok = false
	for _, sev := range []string{INFO, WARN, ERROR, STAR, HASH} {
		if parts[1] == sev {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}

	return true
}

func randSev() string {
	return []string{INFO, WARN, ERROR}[rand.Intn(3)]
}
