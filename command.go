package main

type Command struct {
	Command string
	Args    string
}

type CommandHandler func(client Bot, command, args string) error
