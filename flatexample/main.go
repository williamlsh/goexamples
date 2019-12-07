// Reference: https://medium.com/quick-code/write-a-web-service-with-go-plug-ins-c0472e0645e6?
package main

import (
	"log"
	"plugin"
)

func main() {
	// open plugin file
	plug, err := plugin.Open("plugins/first.so")
	if err != nil {
		log.Fatal(err)
	}
	
	// search for an exported symbol
	symbol, err := plug.Lookup("Talk")
	if err != nil {
		log.Fatal(err)
	}
	
	// call the function
	symbol.(func())()
}
