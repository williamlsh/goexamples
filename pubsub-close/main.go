// Reference: https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ps := NewPubsub()
	ch1 := ps.Subscribe("tech")
	ch2 := ps.Subscribe("travel")
	ch3 := ps.Subscribe("travel")

	listener := func(name string, ch chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
		fmt.Printf("[%s] done\n", name)
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
	ps.Close()
	time.Sleep(50 * time.Millisecond)
}

type Pubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan string
	closed bool
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		subs: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(topic string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Publish(topic, msg string) {
	if ps.closed {
		return
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, sub := range ps.subs {
			for _, ch := range sub {
				close(ch)
			}
		}
	}
}
