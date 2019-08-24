package main

import (
	"fmt"
	"os"
	"time"
)

var (
	tickerForWatcher *time.Ticker
)

const (
	PROCESS_WATCHER_INTERVAL time.Duration = 10 * time.Second
)

func GetProcessStatus(processId int) (string, error) {
	process, err := os.FindProcess(processId)

	if err == nil {
		fmt.Print(true)
		fmt.Printf("Process %d is found", process.Pid)
	} else {
		fmt.Println(false)
	}

	fmt.Println(process)

	return "Watching", nil
}

// This function watches the processId
// Callback function is called when status changes.
// Supports 2 states:
// Exited: Process was killed
// Minimized: The window of the process has been minimized
func JustKeepWatching(processId int) {

	// 1. Start ticker. And just keep watching the process
	tickerForWatcher = time.NewTicker(PROCESS_WATCHER_INTERVAL)
	fmt.Println("Watcher started.")
	go func() {
		for _ = range tickerForWatcher.C {
			_, err := GetProcessStatus(processId)

			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("still looking", processId)
			}
		}
	}()

}
