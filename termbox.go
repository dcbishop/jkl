package main

import "github.com/nsf/termbox-go"

// TermboxDriver is a ConsoleDriver that uses termbox-go.
type TermboxDriver struct {
	events chan Event
	quit   chan interface{}
	width  int
	height int
}

// NewTermboxDriver constructs a new TermboxDriver.
func NewTermboxDriver() TermboxDriver {
	tbd := TermboxDriver{
		width:  0,
		height: 0,
	}

	tbd.initializeChannels()

	termbox.Init()
	tbd.setSize(termbox.Size())

	return tbd
}

func (tbd *TermboxDriver) initializeChannels() {
	if tbd.events == nil {
		tbd.events = make(chan Event)
	}
	if tbd.quit == nil {
		tbd.quit = make(chan interface{})
	}
}

// Size returns the current size of the terminal from Termbox.
func (tbd *TermboxDriver) Size() (width int, height int) {
	return tbd.width, tbd.height
}

// Init initilizes the Termbox library.
func (tbd *TermboxDriver) Init() {
	go tbd.handleEvents()
}

// Close cleansup the Termbox library.
// [TODO]: If this is called on one instance of multiple TermboxDriver
// instances then they all die... Need a refrence count. - 2014-09-30 02:25pm
func (tbd *TermboxDriver) Close() {
	close(tbd.quit)
	tbd.quit = nil
	tbd.initializeChannels()

	termbox.Close()
}

// SetCell sets a character in the console
func (tbd *TermboxDriver) SetCell(x, y int, r rune, fg, bg Color) {
	termbox.SetCell(x, y, r, colorToAttribute(fg), colorToAttribute(bg))
}

// SetCursor sets the Termbox cursor position
func (tbd *TermboxDriver) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

// Events returns a channel of events
func (tbd *TermboxDriver) Events() chan Event {
	return tbd.events
}

// AfterDraw executes a Termbox Flush
func (tbd *TermboxDriver) AfterDraw() {
	termbox.Flush()
}

func colorToAttribute(color Color) termbox.Attribute {
	return color.(termbox.Attribute)
}

func (tbd *TermboxDriver) handleEvents() {
loop:
	for {
		select {
		case <-tbd.quit:
			break loop
		default:
			tbd.handleEvent()
		}
	}
}

func (tbd *TermboxDriver) handleEvent() {
	event := termbox.PollEvent()

	if event.Type == termbox.EventResize {
		tbd.setSize(event.Width, event.Height)
	}

	tbd.events <- termboxEventToInternal(event)
}

func (tbd *TermboxDriver) setSize(width, height int) {
	tbd.width = width
	tbd.height = height
}

func termboxEventToInternal(event termbox.Event) Event {
	return Event{event}
}
