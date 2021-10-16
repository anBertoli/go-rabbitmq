package main

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type task struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func (t task) String() string {
	return fmt.Sprintf("{ name: '%s', level: %d }", t.Name, t.Level)
}

func parseTaskMessage(message amqp.Delivery) (t task, err error) {
	err = json.Unmarshal(message.Body, &t)
	return
}
