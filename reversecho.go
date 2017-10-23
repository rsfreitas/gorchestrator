//The reverse echo bot
package main

import (
	"github.com/processone/gox/xmpp"
)

type EchoBot struct {
	name string
}

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

//
// The BotModel API implementation
//

func (b *EchoBot) HandleMessage(message xmpp.ClientMessage, session ChatSession) []ChatMessage {
	var answers []ChatMessage

	//prepare our answer
	answers = append(answers, ChatMessage{Message: reverse(message.Body)})
	return answers
}

func (b *EchoBot) Name() string {
	return b.name
}

//
// The required method to register ourselves inside the supported bots
//

func NewEchoBot() (BotModel, error) {
	return &EchoBot{
		name: "echo-bot",
	}, nil
}
