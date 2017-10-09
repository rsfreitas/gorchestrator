//The interface to implement a chatbot
package main

import (
	"github.com/processone/gox/xmpp"
)

// TODO: Convert the xmpp.ClientMessage to an internal type

type BotModel interface {
	Name() string
	HandleMessage(message xmpp.ClientMessage) string
}
