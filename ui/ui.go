package ui

import (
	"log"

	"github.com/nsf/termbox-go"
)

// TermboxUI handles ui using termbox-go
type TermboxUI struct {
	quit chan bool
}

// Run enters the main ui loop untill close is called on the quit channel.
func (tbw *TermboxUI) Run() {
	tbw.initializeQuitChannel()

	termbox.Init()

	go tbw.handleEvents()
	tbw.waitForQuit()
}

func (tbw *TermboxUI) initializeQuitChannel() {
	tbw.quit = make(chan bool)
}

// Stop terminates the Run loop.
func (tbw *TermboxUI) Stop() {
	close(tbw.quit)
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
	log.Println(event)
}

func (tbw *TermboxUI) waitForQuit() {
	select {
	case <-tbw.quit:
		termbox.Close()
		return
	}
}
