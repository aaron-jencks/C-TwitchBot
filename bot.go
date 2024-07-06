package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aaronjencks/gitchbot/storage"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

type Bot interface {
	Channel() string
	Join(channel string) error
	Depart(channel string) error
	Say(message string) error
	Whisper(user, message string) error
	Storage() storage.StorageBacking
	RegisterHandler(name string, handler CommandHandler)
	HandlerExists(name string) bool
	UnregisterHandler(name string)
	Loop()
}

type BasicTwitchBot struct {
	username string
	channel  string
	client   *twitch.Client
	handlers map[string]CommandHandler
	storage  storage.StorageBacking
}

func CreateBasicTwitchBot(username, oauth string, backer storage.StorageBacking) *BasicTwitchBot {
	result := BasicTwitchBot{
		username: username,
		client:   twitch.NewClient(username, oauth),
		handlers: map[string]CommandHandler{},
		storage:  backer,
	}
	return &result
}

func (bb *BasicTwitchBot) Channel() string {
	return bb.channel
}

func (bb *BasicTwitchBot) Storage() storage.StorageBacking {
	return bb.storage
}

func (bb *BasicTwitchBot) HandlerExists(name string) bool {
	_, ok := bb.handlers[name]
	return ok
}

func (bb *BasicTwitchBot) RegisterHandler(name string, handler CommandHandler) {
	bb.handlers[name] = handler
}

func (bb *BasicTwitchBot) UnregisterHandler(name string) {
	delete(bb.handlers, name)
}

func (bb *BasicTwitchBot) Join(channel string) error {
	bb.client.Join(channel)
	bb.channel = channel
	return nil
}

func (bb *BasicTwitchBot) Depart(channel string) error {
	bb.client.Depart(channel)
	bb.channel = ""
	return nil
}

func (bb *BasicTwitchBot) Say(message string) error {
	bb.client.Say(bb.channel, message)
	return nil
}

func (bb *BasicTwitchBot) Whisper(user, message string) error {
	bb.client.Say(bb.channel, fmt.Sprintf("/w %s %s", user, message))
	return nil
}

func (bb *BasicTwitchBot) Loop() {
	bb.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf("%s: %s\n", message.User.DisplayName, message.Message)

		if message.User.DisplayName != bb.username {
			TimerMarkMessageReceived() // this helps avoid spam
		}

		if ContainsCommand(message.Message) {
			cmd, err := ParseCommand(message.Message)
			if err != nil {
				log.Printf("failed to parse command message: %v\n", err)
				return
			}
			handler, ok := bb.handlers[cmd.Command]
			if !ok {
				log.Printf("no handler found for command \"%s\"\n", cmd.Command)
				return
			}
			err = handler(bb, ReducedMessage{
				User:    message.User,
				Channel: message.Channel,
				Message: message.Message,
			}, cmd)
			if err != nil {
				log.Printf("failed to handle command \"%s\" with params: \"%s\": %v\n", cmd.Command, cmd.Args, err)
			}
		}
	})

	go func() {
		timer := time.NewTicker(time.Second)
		for {
			<-timer.C
			err := HandleTimers(bb)
			if err != nil {
				log.Printf("failed to handle timers: %v\n", err)
				continue
			}
		}
	}()

	log.Printf("Bot %s started...\n", bb.username)
	err := bb.client.Connect()
	if err != nil {
		log.Println(err)
	}
}
