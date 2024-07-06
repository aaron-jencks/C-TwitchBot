package main

import (
	twitch "github.com/gempir/go-twitch-irc/v4"
	"github.com/oriser/regroup"
)

type Command struct {
	Command string
	Args    string
}

var CMD_REGEX = regroup.MustCompile(`!(?P<command>\w+)(\s+(?P<args>.+))?`)

func ContainsCommand(line string) bool {
	_, err := CMD_REGEX.Groups(line)
	return err == nil
}

func ParseCommand(line string) (cmd Command, err error) {
	matches, err := CMD_REGEX.Groups(line)
	if err != nil {
		return
	}
	cmd.Command = matches["command"]
	cmd.Args = matches["args"]
	return
}

type ReducedMessage struct {
	User    twitch.User
	Channel string
	Message string
}

func (rm ReducedMessage) IsModerator() bool {
	_, broad := rm.User.Badges["broadcaster"]
	_, mod := rm.User.Badges["moderator"]
	return broad || mod
}

type CommandHandler func(client Bot, msg ReducedMessage, command Command) error
