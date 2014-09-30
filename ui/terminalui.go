package ui

import (
	"time"

	"github.com/dcbishop/jkl/editor"
	"github.com/dcbishop/jkl/runegrid"
	"github.com/dcbishop/jkl/service"
	"github.com/nsf/termbox-go"
)

// Color is red and yellow and pink and green...
type Color interface{}

// ConsoleDriver handles the input/output of the TerminalUI to/from the terminal.
type ConsoleDriver interface {
	Size() (width int, height int)
	Init()
	Close()
	SetCell(x, y int, r rune, fg Color, bg Color)
	Events() chan Event
	SetCursor(x, y int)
	AfterDraw()
}

// TerminalUI a text based user interface renderer.
type TerminalUI struct {
	quit    chan bool
	state   service.State
	console ConsoleDriver
}

// Run enters the main UI loop untill Stop() is called.
func (tui *TerminalUI) Run() {
	if tui.state.SetRunning() != nil {
		panic("UI already running.")
		return
	}
	defer tui.state.SetStopped()

	tui.initialize()
	tui.waitForQuit()
	tui.cleanUp()
}

// Running returns true if ui.Run() was called but ui.Stop() hasn't been.
func (tui *TerminalUI) Running() bool {
	return tui.state.Running()
}

// Stop terminates the Run loop.
func (tui *TerminalUI) Stop() {
	if tui.quit == nil {
		return
	}
	close(tui.quit)
	service.WaitUntilStopped(tui, time.Second)
}

// Events gets the channel that emits events
func (tui *TerminalUI) Events() <-chan Event {
	return tui.console.Events()
}

// Redraw updates the display
func (tui *TerminalUI) Redraw(editor editor.Editor) {
	if !tui.state.Running() {
		return
	}
	width, height := tui.console.Size()
	// [TODO]: Cache runegrid and change on resize only - 2014-09-27 10:10pm

	defer tui.console.AfterDraw()

	grid := runegrid.New(width, height)
	grid.RenderEditor(editor)

	tui.renderGrid(&grid)

	if editor.CurrentPane().Cursor() == nil {
		return
	}

	xPos := editor.CurrentPane().Cursor().XPos()
	yPos := editor.CurrentPane().Cursor().LineNumber()

	if editor.Settings().Borders && editor.Settings().OuterBorder {
		xPos++
		yPos++
	}

	tui.console.SetCursor(xPos, yPos)
}

func (tui *TerminalUI) renderGrid(grid *runegrid.RuneGrid) {
	for y, l := range grid.Cells() {
		for x, r := range l {
			tui.console.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorRed)
		}
	}
}

func (tui *TerminalUI) initialize() {
	tui.initializeConsoleDriver()
	tui.initializeQuitChannel()
	tui.console.Init()
}

func (tui *TerminalUI) cleanUp() {
	tui.console.Close()
}

func (tui *TerminalUI) initializeConsoleDriver() {
	if tui.console == nil {
		tui.console = &TermboxDriver{}
	}
}

func (tui *TerminalUI) initializeQuitChannel() {
	tui.quit = make(chan bool)
}
