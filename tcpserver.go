// +build tcp_server

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(data) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", data)

		now := time.Now().Format(time.RFC3339) + "\n"
		c.Write([]byte(now))
	}
}
