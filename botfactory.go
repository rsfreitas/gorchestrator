//The bot factory
package main

import (
	"errors"
)

type BotFactory func() (BotModel, error)

var botFactories = make(map[string]BotFactory)

func register(name string, factory BotFactory) {
	if factory == nil {
		//panic
	}

	_, registered := botFactories[name]

	if registered {
		//error
	}

	botFactories[name] = factory
}

//registerKnownBots register all supported bots.
func registerKnownBots() {
	register("echo-bot", NewEchoBot)
}

func CreateBot(name string) (BotModel, error) {
	registerKnownBots()
	bot, ok := botFactories[name]

	if !ok {
		return nil, errors.New("")
	}

	return bot()
}
