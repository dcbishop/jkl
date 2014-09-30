package app

import (
	"log"
	"time"

	"github.com/dcbishop/fileaccessor"
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
	if event.Ch == 'j' {
		app.editor.CurrentPane().Cursor().MoveDownLine()
	}
	if event.Ch == 'k' {
		app.editor.CurrentPane().Cursor().MoveUpLine()
	}
}

// Update processes input and redraws the app.
func (app *App) Update() {
	app.UI.Redraw(app.editor)
}
