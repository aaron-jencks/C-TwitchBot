package main

import "fmt"

var counters = map[string]int{}

func CreateCounterHandler(b Bot, name string, initial int, statusPrefix string) {
	b.Storage().CreateCounter(name, initial, statusPrefix)
	b.RegisterHandler(name, func(client Bot, msg ReducedMessage, command Command) error {
		backing := client.Storage()
		current, err := backing.RetrieveCounter(name)
		if err != nil {
			return err
		}
		current++
		err = backing.UpdateCounter(name, current)
		if err != nil {
			return err
		}
		return client.Say(msg.Channel, fmt.Sprintf("%s: %d\n", statusPrefix, counters[name]))
	})
}
