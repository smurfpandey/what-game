package main

import (
	"fmt"
	"os"
)

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func main() {

	//
	// foundGame := false

	JustKeepLooking(func(gotTheGame Game) {
		fmt.Println("Got this: ", gotTheGame)
	})
	for {
	}
}
