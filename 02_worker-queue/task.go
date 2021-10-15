package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type task struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (t task) String() string {
	return fmt.Sprintf("{ name: '%s', level: %d }", t.Name, t.Level)
}

func taskFromArgs(arg string) string {
	var t task
	if arg == "" {
		t = task{"hello", 1}
	} else {
		level := bytes.Count([]byte(arg), []byte("."))
		t = task{arg, level}
	}
	encoded, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return string(encoded)
}
