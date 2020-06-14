package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("* * * * * *", func() {
		fmt.Println("Every second.")
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	<-ch
	c.Stop()
}
