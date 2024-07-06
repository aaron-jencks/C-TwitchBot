package main

import (
	"time"
)

func CreateTimer(b Bot, name, message string, interval time.Duration) error {
	return b.Storage().CreateTimer(name, message, interval)
}

func HandleTimers(b Bot) error {
	backer := b.Storage()
	timers, err := backer.ListTimers()
	if err != nil {
		return nil
	}
	t := time.Now()
	for name, val := range timers {
		if val.After(t) {
			continue
		}
		msg, _, _, err := backer.RetrieveTimer(name)
		if err != nil {
			return err
		}
		err = b.Say(msg)
		if err != nil {
			return err
		}
		err = backer.ResetTimer(name)
		if err != nil {
			return err
		}
	}
	return nil
}
