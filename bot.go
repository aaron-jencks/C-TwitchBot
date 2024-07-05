package main

import (
	"fmt"
	"log"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

type Bot interface {
	Join(channel string) error
	Depart(channel string) error
	Say(channel, message string) error
	Whisper(channel, user, message string) error
	Loop()
}

type BasicTwitchBot struct {
	username string
	client   *twitch.Client
	Handlers map[string]CommandHandler
}

func CreateBasicTwitchBot(username, oauth string) *BasicTwitchBot {
	result := BasicTwitchBot{
		username: username,
		client:   twitch.NewClient(username, oauth),
		Handlers: map[string]CommandHandler{},
	}
	return &result
}

func (bb *BasicTwitchBot) Join(channel string) error {
	bb.client.Join(channel)
	return nil
}

func (bb *BasicTwitchBot) Depart(channel string) error {
	bb.client.Depart(channel)
	return nil
}

func (bb *BasicTwitchBot) Say(channel, message string) error {
	bb.client.Say(channel, message)
	return nil
}

func (bb *BasicTwitchBot) Whisper(channel, user, message string) error {
	bb.client.Say(channel, fmt.Sprintf("/w %s %s", user, message))
	return nil
}

func (bb *BasicTwitchBot) Loop() {
	bb.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf("%s: %s\n", message.User.DisplayName, message.Message)
	})

	log.Printf("Bot %s started...\n", bb.username)
	err := bb.client.Connect()
	if err != nil {
		log.Println(err)
	}
}
