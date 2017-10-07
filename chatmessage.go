//Handles different kind of received messages
package main

import (
	"fmt"
)

type MessageType int

const (
	Message MessageType = iota
	ControlMessage
)

type ChatMessage struct {
	Type MessageType
	From string
	To   string
}

func (p ChatMessage) String() string {
	return fmt.Sprintf("{Type: %d, From: %s, To: %s}",
		p.Type, p.From, p.To)
}
