package main

import (
	ps "github.com/keybase/go-ps"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type Game struct {
	Name      string `toml:"name" json:"name"`
	Exec      string `toml:"exec_name" json:"exec_name"`
	Website   string `toml:"website_url" json:"website_url"`
	Store     string `toml:"store_url" json:"store_url"`
	ProcessId int
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
	GAME_LOOKER_INTERVAL time.Duration = 10 * time.Second
)

func LoadGamesDB() {
	bytData, err := ioutil.ReadFile("game-list.toml")
	if err != nil {
		log.Fatal(err)
	}

	err = toml.Unmarshal(bytData, &gamesDB)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Found " + strconv.Itoa(len(gamesDB.Games)) + " games in list")
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

		// if processExecName is ApplicationFrameHost.exe, that means a UWP app is running
		// we should find the real process details instead
		if processExecName == UWP_HOST_APP {
			procId = GetUWPAppProcess()
			yoProcess, _ = ps.FindProcess(procId)
			processExecName = yoProcess.Executable()
		}

		foundGame, milaKya := IsThisAGame(appName, processExecName)
		if milaKya {
			foundGame.ProcessId = procId
		}

		return foundGame, milaKya
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
	tickerForLooker = time.NewTicker(GAME_LOOKER_INTERVAL)
	go func() {
		for _ = range tickerForLooker.C {
			foundGame, gotIt := FindGame()

			if gotIt {
				callback(foundGame)
				tickerForLooker.Stop()
			}
		}
	}()
}
