package ui

import (
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/nsf/termbox-go"
)

// Service implementations have a Run() loop that does stuff and blocks untill Stop() is called
type Service interface {
	Run()
	Stop()
	Running() bool
}

// WaitUntilServiceRunning will wait until service.Running() returns state,
// will return an error if it took longer than timeout.
func WaitUntilServiceRunning(service Service, state bool, timeout time.Duration) error {
	start := time.Now()

	for !service.Running() == state {
		time.Sleep(time.Millisecond)
		if time.Since(start) >= timeout {
			return errors.New("Time out")
		}
	}

	return nil
}

// UI handles input and displaying the app
type UI interface {
	Service
}

// TermboxUI handles ui using termbox-go
type TermboxUI struct {
	quit    chan bool
	running uint32
}

// Run enters the main UI loop untill Stop() is called.
func (tbw *TermboxUI) Run() {
	if tbw.setRunning(true) {
		panic("UI already running.")
		return
	}
	defer tbw.setRunning(false)

	tbw.initialize()

	go tbw.handleEvents()
	tbw.waitForQuit()
	tbw.cleanUp()
}

func (tbw *TermboxUI) setRunning(state bool) (wasAlreadySet bool) {
	newNum := uint32(0)

	if state {
		newNum = 1
	}

	previousState := atomic.SwapUint32(&tbw.running, newNum)

	if previousState == 0 && !state {
		return true
	}

	if previousState == 1 && state {
		return true
	}

	return false
}

// Running returns true if ui.Run() was called but ui.Stop() hasn't been.
func (tbw *TermboxUI) Running() bool {
	return tbw.running == 1
}

// Stop terminates the Run loop.
func (tbw *TermboxUI) Stop() {
	if tbw.quit == nil {
		return
	}
	close(tbw.quit)
	WaitUntilServiceRunning(tbw, false, time.Second)
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
