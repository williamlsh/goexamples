package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func main() {
	if err := publish(*uri, *exchangeName, *exchangeType, *routingKey, *body, *reliable); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("published %dB OK", len(*body))
}

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {
	log.Printf("dialing %q", amqpURI)
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer conn.Close()

	log.Printf("got connection, getting channel")
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	log.Printf("got channel, declaring %q exchange (%q)", exchangeType, exchange)
	if err := ch.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := ch.Confirm(false); err != nil {
			return fmt.Errorf("channel could not be put into confirm mode: %s", err)
		}

		confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	log.Printf("declared exchange, publishing %dB body (%q)", len(body), body)
	if err := ch.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			// a bunch of application/implementation-specific
		},
	); err != nil {
		return fmt.Errorf("exchange publish: %s", err)
	}

	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delevery with delevery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delevery of delevery tag: %d", confirmed.DeliveryTag)
	}
}
