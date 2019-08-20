package main

import (
	"fmt"
	"os"
	"time"
)

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func main() {

	gamesDB := LoadGamesDB()
	foundGame := false

	for {

		fmt.Println(foundGame)
		time.Sleep(10 * time.Second)
		if foundGame {
			fmt.Println("We already have a game running. Sleeping for 10mins")
			time.Sleep(10 * time.Second)
			continue
		}

		_, foundGame = FindGame(gamesDB)
	}
}
