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

}

// OpenFiles opens a list of files into buffers.
func (app *App) OpenFiles(filenames []string) {
	for _, filename := range filenames {
		app.OpenFile(filename)
	}
}

// OpenFile opens a file.
func (app *App) OpenFile(filename string) {
	buffer := NewBuffer()

	if data, err := app.fa.ReadFile(filename); err == nil {
		buffer.data = data
	} else {
		// [TODO]: Error handling. Non-existant files shouldn't cause a problem,
		// but ones that do exist but can't open should show an error.
		buffer.data = []byte{}
	}

	buffer.filename = filename

	app.AddBuffer(buffer)
}

// AddBuffer adds a buffer to the list of buffers
func (app *App) AddBuffer(buffer Buffer) {
	app.buffers = append(app.buffers, buffer)
}

// NewBuffer constructs a new Buffer object containing data.
func NewBuffer() Buffer {
	return Buffer{}
}
