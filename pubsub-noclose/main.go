// Reference: https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ps := NewPubsub()
	makesub := func(topic string) chan string {
		ch := make(chan string, 1)
		ps.Subscribe(topic, ch)
		return ch
	}
	ch1 := makesub("tech")
	ch2 := makesub("travel")
	ch3 := makesub("travel")

	listener := func(name string, ch chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
	}

	go listener("1", ch1)
	go listener("2", ch2)
	go listener("3", ch3)

	pub := func(topic, msg string) {
		fmt.Printf("Publishing @%s: %s\n", topic, msg)
		ps.Publish(topic, msg)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	pub("tech", "tablets")
	pub("health", "vitamins")
	pub("tech", "robots")
	pub("travel", "beaches")
	pub("travel", "hiking")
	pub("tech", "drones")

	time.Sleep(50 * time.Millisecond)
}

type Pubsub struct {
	mu   sync.RWMutex
	subs map[string][]chan string
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		subs: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(topic string, ch chan string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subs[topic] = append(ps.subs[topic], ch)
}

func (ps *Pubsub) Publish(topic, msg string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}
