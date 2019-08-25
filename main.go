package main

import (
	"fmt"
	"os"
)

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func startLogic() {
	JustKeepLooking(func(gotTheGame Game) {
		fmt.Println("So you are playing", gotTheGame.Name)
		fmt.Println("I am tracking you now", gotTheGame.ProcessId)
		JustKeepWatching(gotTheGame.ProcessId, func(statusChange string) {
			fmt.Println(statusChange, "now what? Let's start again?")
			startLogic()
		})
	})
}

func main() {

	//startLogic()
	for {
	}
}
