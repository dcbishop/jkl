package ui

import (
	"time"

	"github.com/dcbishop/jkl/editor"
	"github.com/dcbishop/jkl/runegrid"
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
	service.WaitUntilStopped(tbw, time.Second)
}

// Events gets the channel that emits events
func (tbw *TermboxUI) Events() <-chan Event {
	return tbw.events
}

// Redraw updates the display
func (tbw *TermboxUI) Redraw(editor editor.Editor) {
	if !tbw.state.Running() {
		return
	}
	width, height := termbox.Size()
	// [TODO]: Cache runegrid and change on resize only - 2014-09-27 10:10pm

	grid := runegrid.New(width, height)
	grid.RenderEditor(editor)

	tbw.renderGrid(&grid)
	termbox.SetCursor(1, 1)
}

func (tbw *TermboxUI) renderGrid(grid *runegrid.RuneGrid) {
	for y, l := range grid.Cells() {
		for x, r := range l {
			termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorRed)
		}
	}
	termbox.Flush()
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
