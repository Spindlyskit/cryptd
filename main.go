package main

import (
	"flag"
	"fmt"
)

func main() {
	const (
		portDefault     = 8080
		portDescription = "The port to listen on"
	)
	var portPtr int
	flag.IntVar(&portPtr, "port", portDefault, portDescription)
	flag.IntVar(&portPtr, "p", portDefault, portDescription+" (shorthand)")
	flag.Parse()

	fmt.Println("cryptd rpc server starting on port", portPtr)
	// TODO Start an rpc server that serves crackers
}
