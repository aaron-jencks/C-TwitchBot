package main

import (
	"sync"
	"time"
)

func CreateTimer(b Bot, name, message string, interval time.Duration) error {
	return b.Storage().CreateTimer(name, message, interval)
}

var TIMER_LOCK sync.Mutex = sync.Mutex{}
var TIMER_MSG_INDICATOR bool = true
var LAST_MESSAGE time.Time = time.Now()

func TimerMarkMessageReceived() {
	TIMER_LOCK.Lock()
	defer TIMER_LOCK.Unlock()
	TIMER_MSG_INDICATOR = true
	LAST_MESSAGE = time.Now()
}

func HandleTimers(b Bot) error {
	TIMER_LOCK.Lock()
	defer TIMER_LOCK.Unlock()
	if !TIMER_MSG_INDICATOR {
		return nil
	}

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

	if time.Since(LAST_MESSAGE) > (15 * time.Minute) {
		TIMER_MSG_INDICATOR = false
	}

	return nil
}
