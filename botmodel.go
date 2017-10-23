//The interface to implement a chatbot
package main

import (
	"github.com/processone/gox/xmpp"
)

// TODO: Convert the xmpp.ClientMessage to an internal type

type BotModel interface {
	//Name must return the chatbot name.
	Name() string

	//HandleMessage must be responsible to manipulate the received message and
	//prepare answers, as a slice of ChatMessage, to be returned.
	HandleMessage(message xmpp.ClientMessage, session ChatSession) []ChatMessage
}
