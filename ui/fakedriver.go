package ui

import "github.com/dcbishop/jkl/runegrid"

// FakeDriver is a ConsoleDriver that doesn't do anything.
type FakeDriver struct {
	EventChan chan Event
	Width     int
	Height    int
	CursorX   int
	CursorY   int
	Grid      runegrid.RuneGrid
}

func NewFakeDriver() FakeDriver {
	fd := FakeDriver{}
	return fd
}

func (fd *FakeDriver) Init() {
	fd.EventChan = make(chan Event)
	if fd.Width == 0 {
		fd.Width = 80
	}
	if fd.Height == 0 {
		fd.Height = 24
	}
	fd.Grid = runegrid.New(fd.Width, fd.Height)
}

func (fd *FakeDriver) Events() chan Event {
	return fd.EventChan
}

func (fd *FakeDriver) Size() (width int, height int) {
	return fd.Width, fd.Height
}

func (fd *FakeDriver) Close() {}
func (fd *FakeDriver) SetCell(x, y int, r rune, fg Color, bg Color) {
	fd.Grid.SetCell(x, y, r)
}
func (fd *FakeDriver) SetCursor(x, y int) {
	fd.CursorX = x
	fd.CursorY = y
}
func (fd *FakeDriver) AfterDraw() {}
