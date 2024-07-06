package main

import (
	"fmt"
	"log"
	"time"

	"github.com/oriser/regroup"
)

type HelpEntry struct {
	Username string
	Message  string
	Code     string
}

type HelpCommand struct {
	Command string
	Entry   HelpEntry
	Error   string
}

var HELP_REGEX = regroup.MustCompile(`(?P<command>about|position|pop|put)(\s+"(?P<message>.{20,120})"(\s+https://pastebin\.com/(?P<url>\w+))?)?`)

const HELP_USAGE = "!help (about|position|pop|put \"message/request (20-120 chars)\" [https://pastebin.com/...])"

var helpQueue []HelpEntry

func parseHelpRequest(username, line string) (entry HelpCommand, err error) {
	matches, err := HELP_REGEX.Groups(line)
	if err != nil {
		return
	}
	entry.Command = matches["command"]
	var ok bool
	switch matches["command"] {
	case "put":
		entry.Entry.Username = username
		entry.Entry.Message, ok = matches["message"]
		if !ok {
			entry.Error = "message is required for enqueueing help requests"
		}
		entry.Entry.Code = matches["url"]
	}
	return
}

func getUserHelpPosition(username string) int {
	for hi, entry := range helpQueue {
		if entry.Username == username {
			return hi
		}
	}
	return -1
}

func CreateProgrammingHelpQueue(b Bot) error {
	log.Println("creating hooks for programming help queue")
	b.RegisterHandler("help", func(client Bot, msg ReducedMessage, command Command) error {
		entry, err := parseHelpRequest(msg.User.DisplayName, msg.Message)
		if err != nil {
			return client.Say(fmt.Sprintf("@%s that usage is incorrect, correct usage is: %s", msg.User.DisplayName, HELP_USAGE))
		}
		if entry.Error != "" {
			return client.Say(fmt.Sprintf("@%s there is an error in your request: \"%s\" (see \"!help about\" for usage details)", msg.User.DisplayName, entry.Error))
		}
		switch entry.Command {
		case "about":
			return client.Say(HELP_USAGE)
		case "position":
			idx := getUserHelpPosition(msg.User.DisplayName)
			if idx < 0 {
				return client.Say(fmt.Sprintf("@%s you do not have a request queued at the moment", msg.User.DisplayName))
			} else {
				return client.Say(fmt.Sprintf("@%s you are position %d in the queue", msg.User.DisplayName, idx))
			}
		case "put":
			idx := getUserHelpPosition(msg.User.DisplayName)
			if idx >= 0 {
				return client.Say(fmt.Sprintf("@%s you already have a help request in the queue, please wait your turn, you are at position %d", msg.User.DisplayName, idx))
			}
			helpQueue = append(helpQueue, entry.Entry)
			return client.Say(fmt.Sprintf("@%s you have been added to the queue, you are at position %d", msg.User.DisplayName, len(helpQueue)-1))
		case "pop":
			if !msg.IsModerator() {
				return client.Say(fmt.Sprintf("@%s you must be a moderator to do that", msg.User.DisplayName))
			}
			if len(helpQueue) == 0 {
				return client.Say(fmt.Sprintf("@%s you're all caught up!", msg.User.DisplayName))
			}
			entry := helpQueue[0]
			helpQueue = helpQueue[1:]
			template := fmt.Sprintf("@%s %s asks, \"%s\"", msg.User.DisplayName, entry.Username, entry.Message)
			if entry.Code != "" {
				template += fmt.Sprintf(" they have provided code: https://pastebin.com/%s", entry.Code)
			}
			return client.Say(template)
		}
		return nil
	})
	return CreateTimer(b, "help_timer", "Want to ask a question? Now you can use the queue! See \"!help about\" for usage", 5*time.Minute)
}
