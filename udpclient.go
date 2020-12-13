// +build udp_client

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	s, err := net.ResolveUDPAddr("udp4", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())

	for {
		fmt.Print(">> ")

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(c, text+"\n")

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		buf := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reply: %s\n", buf[:n])
	}
}
