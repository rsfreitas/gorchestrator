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
	Type MessageType
	From string
	To   string

	Attachment string
}

func (p ChatMessage) String() string {
	return fmt.Sprintf("{Type: %d, From: %s, To: %s}",
		p.Type, p.From, p.To)
}
