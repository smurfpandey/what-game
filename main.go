package main

import (
	"fmt"
	"github.com/gonutz/w32"
	"log"
	"os"
	"os/signal"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) run() error {
	startLogic()
	for {
		select {
		case <-p.exit:
			if tickerForLooker != nil {
				tickerForLooker.Stop()
			}

			if tickerForWatcher != nil {
				tickerForWatcher.Stop()
			}

			return nil
		}
	}
}

func startLogic() {
	JustKeepLooking(func(gotTheGame Game) {
		fmt.Println("So you are playing ", gotTheGame.Name)
		fmt.Println("I am tracking you now ", gotTheGame.ProcessId)

		// call POST endpoint
		NotifyGameStarted(gotTheGame)

		JustKeepWatching(gotTheGame.ProcessId, func(statusChange string) {
			fmt.Println(statusChange, " now what? Let's start again?")

			// call DELETE endpoint
			NotifyGameExited()

			startLogic()
		})
	})
}

func main() {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, w32.SW_HIDE)
		}
	}

	prg := &program{}

	errs := make(chan error, 5)

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	prg.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go prg.run()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
}
