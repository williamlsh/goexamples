package main

import (
	"log"
	"plugin"
)

type MyPlug interface {
	Talk()
}

func main() {
	// open plugin file
	plug, err := plugin.Open("plugins/first.so")
	if err != nil {
		log.Fatal(err)
	}
	
	// search for an exported symbol
	symbol, err := plug.Lookup("MyPlugin")
	if err != nil {
		log.Fatal(err)
	}
	
	var myPlugin MyPlug
	myPlugin, ok := symbol.(MyPlug)
	if !ok {
		log.Println("The module has wrong type.")
	}
	
	myPlugin.Talk()
}
