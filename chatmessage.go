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

//A single message object.
type ChatMessage struct {
	Type       MessageType
	From       string
	To         string
	Message    string
	Attachment bool
}

func (p ChatMessage) String() string {
	return fmt.Sprintf("{Type: %d, From: %s, To: %s, Message: %s}",
		p.Type, p.From, p.To, p.Message)
}
