package main

import (
	ps "github.com/keybase/go-ps"
	"github.com/pelletier/go-toml"
	"io/ioutil"
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

func LoadGamesDB() GameList {
	bytData, err := ioutil.ReadFile("game-list.toml")
	if err != nil {
		exitWithMessage(err.Error())
	}

	gamesDB := GameList{}
	err = toml.Unmarshal(bytData, &gamesDB)
	if err != nil {
		exitWithMessage(err.Error())
	}

	return gamesDB
}

func IsThisAGame(gamesDB GameList, appName string, execName string) (Game, bool) {
	for _, game := range gamesDB.Games {
		if game.Name == appName && game.Exec == execName {
			return game, true
		}
	}

	return Game{}, false
}

func FindGame(gamesDB GameList) (Game, bool) {
	if hwnd := GetWindow("GetForegroundWindow"); hwnd != 0 {
		appName := GetWindowText(HWND(hwnd))
		procId := GetWindowProcess(HWND(hwnd))

		yoProcess, _ := ps.FindProcess(procId)
		processExecName := yoProcess.Executable()

		return IsThisAGame(gamesDB, appName, processExecName)
	}

	return Game{}, false
}
