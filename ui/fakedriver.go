package ui

import "github.com/dcbishop/jkl/runegrid"

// FakeDriver is a ConsoleDriver that doesn't do anything.
type FakeDriver struct {
	EventChan chan Event
	width     int
	height    int
	CursorX   int
	CursorY   int
	Grid      runegrid.RuneGrid
}

func NewFakeDriver() FakeDriver {
	fd := FakeDriver{
		EventChan: make(chan Event),
	}
	fd.SetSize(80, 24)
	return fd
}

func (fd *FakeDriver) SetSize(width, height int) {
	fd.width = width
	fd.height = height
	fd.Grid = runegrid.New(fd.width, fd.height)
}

func (fd *FakeDriver) Init() {
}

func (fd *FakeDriver) Events() chan Event {
	return fd.EventChan
}

func (fd *FakeDriver) Size() (width int, height int) {
	return fd.Grid.Size()
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
