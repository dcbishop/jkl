package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dcbishop/jkl/service"
	"github.com/nsf/termbox-go"
	"github.com/spf13/afero"
)

// App is the main program.
type App struct {
	quit   chan interface{}
	UI     UI
	state  service.State
	editor *Editor

	Out    io.Writer
	ErrOut io.Writer
}

// NewApp constructs a new app from the given options.
func NewApp(fs afero.Fs) App {
	editor := NewEditor(fs)
	app := App{
		editor: &editor,
		UI:     nil,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	app.initializeQuitChannel()

	return app
}

// Editor returns the app's editor.
func (app *App) Editor() *Editor {
	return app.editor
}

// SetUI sets the UI to render the app to.
func (app *App) SetUI(ui UI) {
	if app.UI != nil {
		app.UI.Stop()
		service.WaitUntilStopped(app.UI, time.Second)
	}

	app.UI = ui

	if app.UI != nil && app.Running() {
		go app.UI.Run()
	}
}

// SetOut sets the default output stream of the App.
func (app *App) SetOut(out io.Writer) {
	app.Out = out
}

// SetErrOut sets the error output stream of the App.
func (app *App) SetErrOut(eout io.Writer) {
	app.ErrOut = eout
}

// LoadOptions loads the given options.
func (app *App) LoadOptions(options ...Option) {
	for _, o := range options {
		o(app)
	}
}

// Run starts the main loop of the app. Will block until finished.
func (app *App) Run() {
	if app.state.SetRunning() != nil {
		fmt.Fprintf(app.ErrOut, "App already running.")
		return
	}

	if app.UI != nil && !app.UI.Running() {
		go app.UI.Run()
		if service.WaitUntilRunning(app.UI, time.Second) != nil {
			fmt.Fprintf(app.ErrOut, "Could not start UI service.")
		}
	}

	app.loopUntilQuit()
	app.UI.Stop()
	app.state.SetStopped()
}

// Stop shuts everything down and terminates Run(). Blocks untill clean shutdown.
func (app *App) Stop() {
	if !app.Running() {
		return
	}

	app.quit <- true

	if service.WaitUntilStopped(app.UI, time.Second) != nil {
		fmt.Fprintln(app.ErrOut, "UI service did not stop in under a second.")
	}

	if service.WaitUntilStopped(app, time.Second) != nil {
		fmt.Fprintln(app.ErrOut, "App service did not stop in under a second.")
	}
}

// Running returns true if Run() was called but Stop() hasn't been.
func (app *App) Running() bool {
	return app.state.Running()
}

func (app *App) initializeQuitChannel() {
	if app.quit != nil {
		return
	}
	app.quit = make(chan interface{})
}

func (app *App) loopUntilQuit() {
loop:
	for {
		var events <-chan Event
		if app.UI != nil {
			events = app.UI.Events()
		}
		select {
		case <-app.quit:
			break loop
		case event := <-events:
			app.handleEvent(event)
		default:
			app.Update()
		}
	}
}

func (app *App) handleEvent(event Event) {
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

	cursor := app.editor.CurrentPane().Cursor()
	if cursor == nil {
		return
	}

	if event.Ch == 'j' {
		cursor.Move(cursor.DownLine())
	}
	if event.Ch == 'k' {
		cursor.Move(cursor.UpLine())
	}
	if event.Ch == 'h' {
		cursor.Move(cursor.BackCharacter())
	}
	if event.Ch == 'l' {
		cursor.Move(cursor.ForwardCharacter())
	}
}

// Update processes input and redraws the app.
func (app *App) Update() {
	if app.UI != nil {
		app.UI.Redraw(app.editor)
	}
}
