package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	go func() {
		if err := receive(); err != nil {
			log.Fatal(err)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := send(); err != nil {
			log.Fatal(err)
		}
	}

	select {}
}

func send() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err = ch.Confirm(false); err != nil {
		return err
	}

	confirmationCh := make(chan amqp.Confirmation, 1)
	ch.NotifyPublish(confirmationCh)
	go func() {
		for c := range confirmationCh {
			if !c.Ack {
				fmt.Println("Server missed a published message.")
			} else {
				fmt.Println("Sent a message to server.")
			}
		}
	}()

	if err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		return err
	}

	body := "Hello, world!"
	return ch.Publish(
		"logs_direct", // exchange
		"info",        // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{ // msg
			Timestamp:   time.Now(),
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}

func receive() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	if err = ch.ExchangeDeclare(
		"logs_direct",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	if err = ch.QueueBind(
		q.Name,        // queue name
		"info",        // routing key
		"logs_direct", // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		d.Ack(false)
	}

	return nil
}
