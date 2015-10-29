package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/nejstastnejsistene/frotz-slack-bot/rtm"
)

var (
	games = make(map[string]*Zork)
	mutex sync.RWMutex
)

func onMessage(msg rtm.Message, respond chan rtm.Message) {
	if !okMessage(msg) {
		return
	}
	user := "lindenlab"
	text := msg["text"].(string)
	
	if(strings.Contains(text, "botsnack")) {
		defer func() {
			respond <- rtm.NewResponse(msg, "_You are eaten by a grue._")
		}()
		return
	}

	channel, ok := m["channel"].(string)
	
	if(channel != "C0D1YQ44R"){
		return
	}

	text = text[1:len(text)]

	var response string
	defer func() {
		respond <- rtm.NewResponse(msg, "```"+response+"```")
	}()

	mutex.RLock()
	z := games[user]
	mutex.RUnlock()

	if z == nil {
		z, output, err := StartZork("dfrotz", "ZORK1.DAT")
		if err != nil {
			response = fmt.Sprintf("[error: %s]", err)
		} else {

			mutex.Lock()
			games[user] = z
			mutex.Unlock()

			response = output
		}
	} else {
		output, err := z.ExecuteCommand(text)
		if err != nil {
			mutex.Lock()
			delete(games, user)
			mutex.Unlock()

			response = fmt.Sprintf("[error: %s]", err)
			if err == CleanExit {
				response = "[process exited cleanly]"
			}
		} else {
			response = output
		}
	}

	channel := msg["channel"].(string)
	log.Printf("%s: > %s", channel, text)
	for _, line := range strings.Split(response, "\n") {
		log.Printf("%s: %s\n", channel, line)
	}

}

func okMessage(m rtm.Message) bool {
	// Ignore reply_to messages; these are for already sent messages.
	if _, ok := m["reply_to"]; ok {
		return false
	}
	msgType, ok := m["type"].(string)
	if !ok {
		return false
	}
	channel, ok := m["channel"].(string)
	if !ok {
		return false
	}
	if _, ok := m["user"].(string); !ok {
		return false
	}
	text, ok := m["text"].(string)
	if !ok {
		return false
	}

	log.Printf(channel)

	return msgType == "message" && len(channel) > 0 && len(text) > 0 && text[0] == 'z'
}

func main() {
	var token string
	if token = os.Getenv("TOKEN"); token == "" {
		log.Fatal("TOKEN not specified")
	}
	rtm.LoopForever(token, onMessage)
}
