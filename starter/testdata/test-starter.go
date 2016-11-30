package main

import "github.com/tredoe/goutil/starter"

func main() {
	defer starter.ExitStatus() // First defer

	go func() {
		// Code to handle the service conections.
	}()

	starter.Wait(nil, false)
}
