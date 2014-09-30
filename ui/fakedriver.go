package ui

// FakeDriver is a ConsoleDriver that doesn't do anything.
type FakeDriver struct {
	EventChan chan Event
	Width     int
	Height    int
}

func (fd *FakeDriver) Init() {
	fd.EventChan = make(chan Event)
	fd.Width = 80
	fd.Height = 24
}

func (fd *FakeDriver) Events() chan Event {
	return fd.EventChan
}

func (fd *FakeDriver) Size() (width int, height int) {
	return fd.Width, fd.Height
}

func (fd *FakeDriver) Close()                                       {}
func (fd *FakeDriver) SetCell(x, y int, r rune, fg Color, bg Color) {}
func (fd *FakeDriver) SetCursor(x, y int)                           {}
func (fd *FakeDriver) AfterDraw()                                   {}
