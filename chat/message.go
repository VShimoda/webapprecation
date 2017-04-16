package main

import (
	"time"
)

// message is one message
type message struct {
	Name    string
	Message string
	When    time.Time
}
