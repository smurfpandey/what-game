package main

import (
	"fmt"
	ps "github.com/keybase/go-ps"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"time"
)

type Game struct {
	Name    string `toml:"name"`
	Exec    string `toml:"exec_name"`
	Website string `toml:"website_url"`
	Store   string `toml:"store_url"`
}

type GameList struct {
	Games []Game `toml:"games"`
}

type CallbackWhenFound func(foundGame Game)

var (
	tickerForLooker *time.Ticker
	gamesDB         GameList
)

const (
	LOOKER_INTERVAL time.Duration = 10 * time.Second
)

func LoadGamesDB() {
	bytData, err := ioutil.ReadFile("game-list.toml")
	if err != nil {
		exitWithMessage(err.Error())
	}

	err = toml.Unmarshal(bytData, &gamesDB)
	if err != nil {
		exitWithMessage(err.Error())
	}

	fmt.Printf("Loaded %d games.\n", len(gamesDB.Games))
}

func IsThisAGame(appName string, execName string) (Game, bool) {
	for _, game := range gamesDB.Games {
		if game.Name == appName && game.Exec == execName {
			return game, true
		}
	}

	return Game{}, false
}

func FindGame() (Game, bool) {
	if hwnd := GetWindow("GetForegroundWindow"); hwnd != 0 {
		appName := GetWindowText(HWND(hwnd))
		procId := GetWindowProcess(HWND(hwnd))

		yoProcess, _ := ps.FindProcess(procId)
		processExecName := yoProcess.Executable()

		return IsThisAGame(appName, processExecName)
	}

	return Game{}, false
}

// JustKeepLooking once initiated, starts a Ticker to continously look for a active Game window
// Ticker runs every 10sec till a game is found.
// Once found, Ticker is stopped, and callback function is invoked
func JustKeepLooking(callback CallbackWhenFound) {

	// 1. Load gamesDB from file
	LoadGamesDB()

	// 2. Start ticker. And just keep looking
	tickerForLooker = time.NewTicker(LOOKER_INTERVAL)
	fmt.Println("Looker started.")
	go func() {
		for t := range tickerForLooker.C {
			foundGame, gotIt := FindGame()

			if gotIt {
				fmt.Println("Got it", t)
				callback(foundGame)
				tickerForLooker.Stop()
			} else {
				fmt.Println("still looking", t)
			}

		}
	}()
}
