package main

import (
	"flag"
	"log"

	"github.com/kardianos/service"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

var logger service.Logger

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
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
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func startLogic() {
	JustKeepLooking(func(gotTheGame Game) {
		logger.Info("So you are playing ", gotTheGame.Name)
		logger.Info("I am tracking you now ", gotTheGame.ProcessId)

		// call POST endpoint
		NotifyGameStarted(gotTheGame)

		JustKeepWatching(gotTheGame.ProcessId, func(statusChange string) {
			logger.Info(statusChange, " now what? Let's start again?")

			// call DELETE endpoint
			NotifyGameExited()

			startLogic()
		})
	})
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "WhatGameYouPlaying",
		DisplayName: "What Game",
		Description: "This service tracks which I am playing",
		Option:      options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
