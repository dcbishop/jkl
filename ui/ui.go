package ui

import (
	"log"

	"github.com/nsf/termbox-go"
)

// TermboxUI handles ui using termbox-go
type TermboxUI struct {
	quit chan bool
}

// Run enters the main UI loop untill Stop() is called.
func (tbw *TermboxUI) Run() {
	tbw.initialize()

	go tbw.handleEvents()
	tbw.waitForQuit()
	tbw.cleanUp()
}

// Stop terminates the Run loop.
func (tbw *TermboxUI) Stop() {
	if tbw.quit == nil {
		return
	}
	close(tbw.quit)
}

func (tbw *TermboxUI) initialize() {
	tbw.initializeQuitChannel()
	termbox.Init()
}

func (tbw *TermboxUI) cleanUp() {
	termbox.Close()
}

func (tbw *TermboxUI) initializeQuitChannel() {
	tbw.quit = make(chan bool)
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
		return
	}
}
