// +build udp_server

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	s, err := net.ResolveUDPAddr("udp4", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	c, err := net.ListenUDP("udp4", s)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	buf := make([]byte, 1024)
	for {
		n, s, err := c.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("-> %s", buf[:n])

		if strings.TrimSpace(string(buf[:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		_, err = c.WriteToUDP([]byte(strconv.Itoa(random(1, 1001))), s)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
