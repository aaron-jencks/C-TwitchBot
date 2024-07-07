package main

import (
	"fmt"
	"log"
	"strings"
)

func generateMappingHandler(name string) CommandHandler {
	return func(client Bot, msg ReducedMessage, command Command) error {
		backing := client.Storage()
		mout, err := backing.RetrieveMapping(name)
		if err != nil {
			return err
		}
		mout = strings.ReplaceAll(mout, "{user}", msg.User.DisplayName)
		return client.Say(mout)
	}
}

func CreateMappingHandler(b Bot, name, message string) error {
	if b.HandlerExists(name) {
		return fmt.Errorf("failed to create mapping for %s, handler already exists", name)
	}
	b.Storage().CreateMapping(name, message)
	b.RegisterHandler(name, generateMappingHandler(name))
	log.Printf("created new mapping handler for %s\n", name)
	return nil
}

func LoadMappingHandlers(b Bot) error {
	mappings, err := b.Storage().ListMappings()
	if err != nil {
		return err
	}
	for name := range mappings {
		b.RegisterHandler(name, generateMappingHandler(name))
		log.Printf("loaded mapping handler for %s\n", name)
	}
	return nil
}
