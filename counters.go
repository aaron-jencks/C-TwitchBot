package main

import (
	"fmt"
	"log"
)

func generateCounterHandler(name string) CommandHandler {
	return func(client Bot, msg ReducedMessage, command Command) error {
		backing := client.Storage()
		current, prefix, err := backing.RetrieveCounter(name)
		if err != nil {
			return err
		}
		current++
		err = backing.UpdateCounter(name, current)
		if err != nil {
			return err
		}
		return client.Say(fmt.Sprintf("%s: %d\n", prefix, current))
	}
}

func CreateCounterHandler(b Bot, name string, initial int, statusPrefix string) error {
	if b.HandlerExists(name) {
		return fmt.Errorf("failed to create counter %s, handler already exists", name)
	}
	b.Storage().CreateCounter(name, initial, statusPrefix)
	b.RegisterHandler(name, generateCounterHandler(name))
	log.Printf("created new counter handler for %s\n", name)
	return nil
}

func LoadCounterHandlers(b Bot) error {
	counters, err := b.Storage().ListCounters()
	if err != nil {
		return err
	}
	for _, counter := range counters {
		b.RegisterHandler(counter, generateCounterHandler(counter))
		log.Printf("loaded counter handler for %s\n", counter)
	}
	return nil
}
