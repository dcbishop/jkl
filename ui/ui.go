package ui

import (
	"log"
	"time"

	"github.com/dcbishop/gim/service"
	"github.com/nsf/termbox-go"
)

// UI handles input and displaying the app
type UI interface {
	service.Service
}

// TermboxUI handles ui using termbox-go
type TermboxUI struct {
	quit  chan bool
	state service.State
}

// Run enters the main UI loop untill Stop() is called.
func (tbw *TermboxUI) Run() {
	if tbw.state.SetRunning() != nil {
		panic("UI already running.")
		return
	}
	defer tbw.state.SetStopped()

	tbw.initialize()

	go tbw.handleEvents()
	tbw.waitForQuit()
	tbw.cleanUp()
}

// Running returns true if ui.Run() was called but ui.Stop() hasn't been.
func (tbw *TermboxUI) Running() bool {
	return tbw.state.Running()
}

// Stop terminates the Run loop.
func (tbw *TermboxUI) Stop() {
	if tbw.quit == nil {
		return
	}
	close(tbw.quit)
	service.WaitUntilRunning(tbw, time.Second)
}

func (tbw *TermboxUI) initialize() {
	tbw.initializeQuitChannel()
	termbox.Init()
}

func (tbw *TermboxUI) cleanUp() {
	termbox.Close()
}

func (tbw *TermboxUI) initializeQuitChannel() {
	tbw.quit = make(chan bool)
}

func (tbw *TermboxUI) handleEvents() {
loop:
	for {
		select {
		case <-tbw.quit:
			break loop
		default:
			tbw.handleEvent()
		}
	}
}

func (tbw *TermboxUI) handleEvent() {
	event := termbox.PollEvent()
	log.Println(event)
}

func (tbw *TermboxUI) waitForQuit() {
	select {
	case <-tbw.quit:
		return
	}
}
