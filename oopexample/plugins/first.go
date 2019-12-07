package main

import "fmt"

type myPlugin string

func (h myPlugin) Talk() {
	fmt.Println("Hello from go plugin!")
}

var MyPlugin myPlugin
