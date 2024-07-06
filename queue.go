package main

import "github.com/oriser/regroup"

type HelpEntry struct {
	Username string
	Message string
	Code string
}

var HELP_REGEX = regroup.MustCompile(`"(?P<message>.{20,120})"\s+(https://pastebin\.com/)?(?P<url>\w+)`)

var helpQueue []HelpEntry

func parseHelpRequest(username, line string) (entry HelpEntry, err error) {
	matches, err := HELP_REGEX.Groups(line)
	if err != nil {
		return
	}
	entry.Username = username
	entry.Message = matches["message"]
	entry.Code = matches["url"]
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

func CreateProgrammingHelpQueue(b Bot) (CommandHandler, error) {
	b.RegisterHandler("help", func(client Bot, msg ReducedMessage, command Command) error {
		idx := getUserHelpPosition(msg.User.DisplayName)
		if idx >= 0 {
			client.Say("@%s you already have a help request in the queue, please wait your turn, you are at position %d", msg.User.DisplayName, idx)
			return nil
		}

	})
}
