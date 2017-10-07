//Functions to handle a chat session with ourselves.
package main

import (
	"fmt"
	"time"
)

type ChatSession struct {
	lastActivity time.Time
	Name         string
	ChatMessage
}

func (p ChatSession) String() string {
	return fmt.Sprintf("{lastActivity: %d, Name: %s, ChatMessage: %v}",
		p.lastActivity, p.Name, p.ChatMessage)
}

//IsActive returns if the current session is active with us or not.
//To be active a session must have received a message before the TTL
//expires.
//
//TTL is given in minutes.
func (p ChatSession) IsActive(ttl float64) bool {
	if time.Now().Sub(p.lastActivity).Minutes() > ttl {
		return false
	}

	return true
}

func DeleteInactives(sessions map[string]ChatSession, ttl float64) {
	for key, value := range sessions {
		if value.IsActive(ttl) == false {
			fmt.Println("Removing session from user", value.Name)
			delete(sessions, key)
		}
	}
}
