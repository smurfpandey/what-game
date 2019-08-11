package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"golang.org/x/sys/windows"
	ps "github.com/keybase/go-ps"
	"os"
	"syscall"
	"time"
	"io/ioutil"
	"unsafe"
)

var (
	mod                          = windows.NewLazyDLL("user32.dll")
	procGetWindowText            = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength      = mod.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = mod.NewProc("GetWindowThreadProcessId")
)

type (
	HANDLE uintptr
	HWND   HANDLE
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

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func GetWindowProcess(hwnd HWND) int {
	procId := 0;
	procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&procId)))

	return procId
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
		if (game.Name == appName && game.Exec == execName) {
			return game, true
		}		
	}

	return Game{}, false
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

		if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
			appName := GetWindowText(HWND(hwnd))
			procId := GetWindowProcess(HWND(hwnd))
			
			yoProcess, _ := ps.FindProcess(procId)
			processExecName := yoProcess.Executable()
			// processPath, _ := yoProcess.Path()

			_, foundGame = IsThisAGame(gamesDB, appName, processExecName)

			if foundGame {				
				fmt.Println("Found a game")
			} else {
				fmt.Println("Not a game")
			}
		}
	}
}
