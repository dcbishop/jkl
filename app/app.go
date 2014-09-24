package app

import (
	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/gim/cli"
)

// App is the main program.
type App struct {
	fa            fileaccessor.FileAccessor
	buffers       []Buffer
	currentBuffer *Buffer
}

// Buffer contains the text to edit
type Buffer struct {
	filename string
	data     []byte
}

// New constructs a new app from the given options.
func New(fa fileaccessor.FileAccessor) App {
	app := App{fa: fa}
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

// NewBuffer constructs a new Buffer object containing data.
func NewBuffer() Buffer {
	return Buffer{}
}
