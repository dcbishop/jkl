package ui

import "github.com/nsf/termbox-go"

// TermboxDriver is a ConsoleDriver that uses termbox-go.
type TermboxDriver struct {
	events chan Event
	quit   chan interface{}
}

// Size returns the current size of the terminal from Termbox.
func (tbd *TermboxDriver) Size() (width int, height int) {
	return termbox.Size()
}

// Init initilizes the Termbox library.
func (tbd *TermboxDriver) Init() {
	tbd.events = make(chan Event)
	tbd.quit = make(chan interface{})
	termbox.Init()
	go tbd.handleEvents()
}

// Close closes the Termbox library.
// [TODO]: If this is called on one instance of multiple TermboxDriver
// instances then they all die... Need a refrence count,
// the entire TermboxDriver  should be a singleton - 2014-09-30 02:25pm
func (tbd *TermboxDriver) Close() {
	close(tbd.quit)
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
	tbd.events <- termboxEventToInternal(event)
}

func termboxEventToInternal(event termbox.Event) Event {
	return Event{event}
}
