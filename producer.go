package main

import (
	"time"
)

// Producer will produce jobs
type Producer struct {
	callbackFunc func(event int)
}

func (p *Producer) start() {
	eventIndex := 0
	for {
		p.callbackFunc(eventIndex)
		eventIndex++
		time.Sleep(time.Millisecond * 500)
	}
}
