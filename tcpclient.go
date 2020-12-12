// +build tcp_client

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
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print(">>")

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(c, text+"\n")

		msg, err := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + msg)

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}
