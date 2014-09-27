package ui

import (
	"time"

	"github.com/dcbishop/jkl/service"
	"github.com/nsf/termbox-go"
)

// TermboxUI handles ui using termbox-go.
type TermboxUI struct {
	quit   chan bool
	events chan Event
	state  service.State
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

// Events gets the channel that emits events
func (tbw *TermboxUI) Events() <-chan Event {
	return tbw.events
}

func (tbw *TermboxUI) initialize() {
	tbw.initializeQuitChannel()
	tbw.initializeEventChannel()
	termbox.Init()
}

func (tbw *TermboxUI) cleanUp() {
	termbox.Close()
}

func (tbw *TermboxUI) initializeQuitChannel() {
	tbw.quit = make(chan bool)
}

func (tbw *TermboxUI) initializeEventChannel() {
	tbw.events = make(chan Event)
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
	tbw.events <- termboxEventToInternal(event)
}

func termboxEventToInternal(event termbox.Event) Event {
	return Event{event}
}

func (tbw *TermboxUI) waitForQuit() {
	select {
	case <-tbw.quit:
		return
	}
}
