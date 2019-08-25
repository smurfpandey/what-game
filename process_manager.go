package main

import (
	ps "github.com/keybase/go-ps"
	"log"
	"time"
)

type CallbackOnStatusChange func(newStatus string)

const (
	PROCESS_STATUS_ACTIVE   = "ACTIVE"
	PROCESS_STATUS_INACTIVE = "INACTIVE" // Minimized
	PROCESS_STATUS_EXITED   = "EXITED"
)

var (
	tickerForWatcher *time.Ticker
)

const (
	PROCESS_WATCHER_INTERVAL time.Duration = 10 * time.Second
)

func GetProcessStatus(processId int) (string, error) {
	yoProcess, err := ps.FindProcess(processId)

	if err == nil {
		if yoProcess == nil {
			return PROCESS_STATUS_EXITED, nil
		} else {
			return PROCESS_STATUS_ACTIVE, nil
		}
	} else {
		return "ERROR", err
	}
}

// This function watches the processId
// Callback function is called when status changes.
// Supports 2 states:
// Exited: Process was killed
// Minimized: The window of the process has been minimized
func JustKeepWatching(processId int, callback CallbackOnStatusChange) {

	// 1. Start ticker. And just keep watching the process
	tickerForWatcher = time.NewTicker(PROCESS_WATCHER_INTERVAL)
	go func() {
		for _ = range tickerForWatcher.C {
			procStatus, err := GetProcessStatus(processId)

			if err != nil {
				log.Fatal(err)
			} else {
				switch procStatus {
				case PROCESS_STATUS_ACTIVE:
					logger.Info("Still playing that game")
				case PROCESS_STATUS_EXITED:
					tickerForWatcher.Stop()
					callback(procStatus)

				}
			}
		}
	}()

}
