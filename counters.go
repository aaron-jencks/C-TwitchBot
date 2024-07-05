package main

import "fmt"

var counters = map[string]int{}

func CreateCounterHandler(name string, initial int, statusPrefix string) CommandHandler {
	counters[name] = initial
	return func(client Bot, msg ReducedMessage, command Command) error {
		counters[name]++
		return client.Say(msg.Channel, fmt.Sprintf("%s: %d\n", statusPrefix, counters[name]))
	}
}
