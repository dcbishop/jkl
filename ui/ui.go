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
	Redraw(editor editor.Interface)
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
	eventChan chan Event
	quit      chan interface{}
}

// NewFakeUI constructs a new FakeUI.
func NewFakeUI() FakeUI {
	fui := FakeUI{}
	fui.initialize()
	return fui
}

// Run starts the service. Will block until Stop() is called.
func (ui *FakeUI) Run() {
	if ui.state.SetRunning() != nil {
		panic("Fake UI already running.")
	}

	ui.initialize()
	ui.loopUntilQuit()
	ui.state.SetStopped()
}

func (ui *FakeUI) initialize() {
	if ui.quit == nil {
		ui.quit = make(chan interface{})
	}
	if ui.eventChan == nil {
		ui.eventChan = make(chan Event)
	}
}

// Stop the service
func (ui *FakeUI) Stop() {
	if !ui.Running() {
		return
	}

	ui.quit <- true
	service.WaitUntilStopped(ui, time.Second)
}

// Running returns the running state of the service
func (ui *FakeUI) Running() bool {
	return ui.state.Running()
}

// Events returns the channel that produces events
func (ui *FakeUI) Events() <-chan Event {
	return ui.eventChan
}

// Redraw updates the display
func (ui *FakeUI) Redraw(editor editor.Interface) {
}

func (ui *FakeUI) loopUntilQuit() {
loop:
	for {
		select {
		case <-ui.quit:
			break loop
		}
	}
}
