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
	JustKeepWatching(18864)
	// JustKeepLooking(func(gotTheGame Game) {
	// 	fmt.Println("So you are playing", gotTheGame.Name)
	// 	fmt.Println("I am tracking you now", gotTheGame.ProcessId)
	// 	JustKeepWatching(gotTheGame.ProcessId)
	// })
	for {
	}
}
