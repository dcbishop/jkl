package ui

import (
	"time"

	"github.com/dcbishop/jkl/editor"
	"github.com/dcbishop/jkl/service"
)

// UI handles input and displaying the app.
type UI interface {
	service.Service
	Input
	Redraw(editor editor.Editor)
}

// Input implementers provide a chanel that generates events.
type Input interface {
	Events() <-chan Event
}

// Event holds information about an event.
type Event struct {
	Data interface{}
}

// FakeUI for disabling output, injecting input and testing.
type FakeUI struct {
	state     service.State
	EventChan chan Event
	Quit      chan bool
}

// Run starts the service. Will block until Stop() is called.
func (ui *FakeUI) Run() {
	if ui.state.SetRunning() != nil {
		panic("Fake UI already running.")
	}

	ui.EventChan = make(chan Event)
	ui.Quit = make(chan bool)

	ui.loopUntilQuit()

	ui.state.SetStopped()
}

// Stop the service
func (ui *FakeUI) Stop() {
	if !ui.Running() {
		return
	}
	close(ui.Quit)
	service.WaitUntilStopped(ui, time.Second)
}

// Running returns the running state of the service
func (ui *FakeUI) Running() bool {
	return ui.state.Running()
}

// Events returns the channel that produces events
func (ui *FakeUI) Events() <-chan Event {
	return ui.EventChan
}

// Redraw updates the display
func (ui *FakeUI) Redraw(editor editor.Editor) {
}

func (ui *FakeUI) loopUntilQuit() {
loop:
	for {
		select {
		case <-ui.Quit:
			break loop
		}
	}
}
