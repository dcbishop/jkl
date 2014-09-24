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

type Buffer struct {
	data []byte
}

// New constructs a new app from the given options.
func New(fa fileaccessor.FileAccessor) App {
	app := App{fa: fa}
	return app
}

// LoadOptions loads the given options.
func (app *App) LoadOptions(options cli.Options) {

}

// OpenFile opens a file.
func (app *App) OpenFile(filename string) {
	data, err := app.fa.ReadFile(filename)
	if err != nil {
		// [TODO]: Error handling...
		return
	}
	buffer := NewBuffer(data)
	app.AddBuffer(buffer)
}

// AddBuffer adds a buffer to the list of buffers
func (app *App) AddBuffer(buffer Buffer) {
	app.buffers = append(app.buffers, buffer)
}

// NewBuffer constructs a new Buffer object containing data.
func NewBuffer(data []byte) Buffer {
	return Buffer{data: data}
}
