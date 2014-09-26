package app

import (
	"fmt"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/gim/cli"
	"github.com/dcbishop/gim/service"
	"github.com/dcbishop/gim/ui"
)

// App is the main program.
type App struct {
	fa            fileaccessor.FileAccessor
	buffers       []Buffer
	currentBuffer *Buffer
	quit          chan interface{}
	ui            ui.UI
	state         service.State
}

// Buffer contains the text to edit
type Buffer struct {
	filename string
	data     []byte
}

// New constructs a new app from the given options.
func New(fa fileaccessor.FileAccessor) App {
	app := App{
		fa: fa,
	}
	return app
}

// LoadOptions loads the given options.
func (app *App) LoadOptions(options cli.Options) {
	app.OpenFiles(options.FilesToOpen)
}

// OpenFiles opens a list of files into buffers and sets the current buffer to the first of the new buffers.
func (app *App) OpenFiles(filenames []string) {
	for i, filename := range filenames {
		buffer := app.openFile(filename)
		bufferPtr := app.AddBuffer(buffer)

		if i == 0 {
			app.currentBuffer = bufferPtr
		}
	}
}

// openFile reads a file, loads it into a new buffer and adds it to the list of buffers
func (app *App) openFile(filename string) Buffer {
	buffer := NewBuffer()

	buffer.filename = filename

	if data, err := app.fa.ReadFile(filename); err == nil {
		buffer.data = data
	} else {
		// [TODO]: Error handling. Non-existant files shouldn't cause a problem,
		// but ones that do exist but can't open should show an error.
		buffer.data = []byte{}
	}

	return buffer
}

// OpenFile opens a file and sets it to the current buffer.
func (app *App) OpenFile(filename string) {
	app.OpenFiles([]string{filename})
}

// AddBuffer adds a buffer to the list of buffers
func (app *App) AddBuffer(buffer Buffer) *Buffer {
	app.buffers = append(app.buffers, buffer)
	return app.LastBuffer()
}

// LastBuffer returns a pointer to the last buffer in the list of buffers
func (app *App) LastBuffer() *Buffer {
	return &app.buffers[len(app.buffers)-1]
}

// SetCurrentBuffer sets the currently visible buffer
func (app *App) SetCurrentBuffer(buffer *Buffer) {
}

// Run starts the main loop of the app. Will block until finished.
func (app *App) Run() {
	if app.state.SetRunning() != nil {
		panic("App already running.")
	}
	app.initialize()

	go app.ui.Run()
	defer app.ui.Stop()

	app.loopUntilQuit()
	app.state.SetStopped()
}

// Stop shuts everything down and terminates Run(). Blocks untill clean shutdown.
func (app *App) Stop() {
	if app.quit == nil {
		return
	}
	close(app.quit)

	if service.WaitUntilStopped(app.ui, time.Second) != nil {
		fmt.Println("UI service did not stop in under a second.")
	}
	if service.WaitUntilStopped(app, time.Second) != nil {
		fmt.Println("App service did not stop in under a second.")
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
	if app.ui == nil {
		app.ui = &ui.TermboxUI{}
	}
}

func (app *App) loopUntilQuit() {
loop:
	for {
		select {
		case <-app.quit:
			break loop
		default:
			app.Update()
		}
	}
}

// Update processes input and redraws the app.
func (app *App) Update() {
}

// NewBuffer constructs a new Buffer object containing data.
func NewBuffer() Buffer {
	return Buffer{}
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
	buffer *Buffer,
	wrap,
	linebrake,
	breakindent bool,
	showbreak string,
	// [TODO]: Should line numbering be done here?
	// Otherwise how do we communicate the line numbers out. - 2014-09-24 11:50pm
) {
	xPos := x
	yPos := y

	for _, r := range buffer.data {
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
