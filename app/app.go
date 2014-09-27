package app

import (
	"log"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/buffer"
	"github.com/dcbishop/jkl/cli"
	"github.com/dcbishop/jkl/editor"
	"github.com/dcbishop/jkl/service"
	"github.com/dcbishop/jkl/ui"
	"github.com/nsf/termbox-go"
)

// App is the main program.
type App struct {
	quit   chan interface{}
	UI     ui.UI
	state  service.State
	editor editor.Editor
}

// New constructs a new app from the given options.
func New(fa fileaccessor.FileAccessor) App {
	editor := editor.New(fa)
	app := App{
		editor: &editor,
	}
	return app
}

// LoadOptions loads the given options.
func (app *App) LoadOptions(options cli.Options) {
	app.editor.OpenFiles(options.FilesToOpen)
}

// Run starts the main loop of the app. Will block until finished.
func (app *App) Run() {
	if app.state.SetRunning() != nil {
		panic("App already running.")
	}
	app.initialize()

	go app.UI.Run()
	defer app.UI.Stop()

	app.loopUntilQuit()
	app.state.SetStopped()
}

// Stop shuts everything down and terminates Run(). Blocks untill clean shutdown.
func (app *App) Stop() {
	if app.quit == nil {
		return
	}

	close(app.quit)

	if service.WaitUntilStopped(app.UI, time.Second) != nil {
		log.Println("UI service did not stop in under a second.")
	}
	if service.WaitUntilStopped(app, time.Second) != nil {
		log.Println("App service did not stop in under a second.")
	}
}

// Running returns true if Run() was called but Stop() hasn't been.
func (app *App) Running() bool {
	return app.state.Running()
}

func (app *App) initialize() {
	app.initializeQuitChannel()
	app.initializeUI()
}

func (app *App) initializeQuitChannel() {
	app.quit = make(chan interface{})
}

func (app *App) initializeUI() {
	if app.UI == nil {
		app.UI = &ui.TermboxUI{}
	}
}

func (app *App) loopUntilQuit() {
loop:
	for {
		select {
		case <-app.quit:
			break loop
		case event := <-app.UI.Events():
			app.handleEvent(event)
		default:
			app.Update()
		}
	}
}

func (app *App) handleEvent(event ui.Event) {
	// [TODO]: Convert all Events to an interal format in the UI layer rather than using termbox directly. - 2014-09-27 11:27am
	switch data := event.Data.(type) {
	case termbox.Event:
		app.handleTermboxEvent(data)
	}
}

func (app *App) handleTermboxEvent(event termbox.Event) {
	switch event.Type {
	case termbox.EventKey:
		app.handleTermboxKeyEvent(event)
	}
}

func (app *App) handleTermboxKeyEvent(event termbox.Event) {
	if event.Ch == 'q' {
		go app.Stop()
	}
}

// Update processes input and redraws the app.
func (app *App) Update() {
	app.UI.Redraw(app.editor)
}

// RuneGrid contains the rendered text UI
type RuneGrid struct {
	width  uint
	height uint
	cells  [][]rune
}

// NewRuneGrid constructs a RuneGrid with the given width and height
func NewRuneGrid(width, height uint) RuneGrid {
	grid := RuneGrid{
		width:  width,
		height: height,
		cells:  make([][]rune, height),
	}

	for i := range grid.cells {
		grid.cells[i] = make([]rune, width)
	}

	return grid
}

// RenderBuffer blits the buffer onto the grid.
// wrap sets line wrapping on
// linebrake sets soft wrapping
func (grid *RuneGrid) RenderBuffer(
	x, y, x2, y2 uint,
	buffer buffer.Buffer,
	wrap,
	linebrake,
	breakindent bool,
	showbreak string,
	// [TODO]: Should line numbering be done here?
	// Otherwise how do we communicate the line numbers out. - 2014-09-24 11:50pm
) {
	xPos := x
	yPos := y

	for _, r := range buffer.Data() {
		if r == '\n' {
			yPos++
			xPos = x
			continue
		}
		if xPos <= x2 && yPos <= y2 {
			grid.SetCell(xPos, yPos, rune(r))
		}
		xPos++
	}
}

// SetCell sets a cell in the RuneGrid to the given rune
func (grid *RuneGrid) SetCell(x, y uint, r rune) {
	if !grid.IsCellValid(x, y) {
		return
	}

	grid.cells[y][x] = r
}

// IsCellValid returns true if the cell coordinates are valid
func (grid *RuneGrid) IsCellValid(x, y uint) bool {
	if x >= grid.width || y >= grid.height {
		return false
	}
	return true
}
