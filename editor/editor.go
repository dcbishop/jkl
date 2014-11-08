package editor

import (
	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/buffer"
)

// Settings stores settings for the editor
// Borders draws pretty borders around panes but takes up some screen space.
// OuterBorder when false turns off just the outer border.
// ShiftWidth is the number of spaces each tab will be displayed as.
// ScrollOffset is the minimum number of lines that will be visible above or below the cursor.
type Settings struct {
	Borders      bool
	OuterBorder  bool
	ShiftWidth   int
	ScrollOffset int
}

// DefaultSettings constructs a default settings.
func DefaultSettings() Settings {
	return Settings{
		Borders:      true,
		OuterBorder:  true,
		ShiftWidth:   4,
		ScrollOffset: 0,
	}
}

// Cursor stores a position in a buffer and handles movement.
type Cursor struct {
	x      int
	line   int
	buffer *buffer.Buffer
}

// Position reutrns the cursors current position.
func (cursor *Cursor) Position() (xPos int, lineNumber int) {
	return cursor.x, cursor.line + 1
}

// Move repositions the cursor at the given coordinates.
func (cursor *Cursor) Move(xPos int, lineNumber int) {
	cursor.x = xPos
	cursor.line = lineNumber - 1
}

// DownLine returns the cursors position one line down.
func (cursor *Cursor) DownLine() (xPos int, lineNumber int) {
	_, err := cursor.buffer.GetLine(cursor.line + 2)
	if err != nil {
		return cursor.Position()
	}
	return cursor.x, cursor.line + 2
}

// UpLine returns the cursors position one line up.
func (cursor *Cursor) UpLine() (xPos int, lineNumber int) {
	if cursor.line == 0 {
		return cursor.Position()
	}
	return cursor.x, cursor.line
}

// BackCharacter returns the cursors position one character back.
func (cursor *Cursor) BackCharacter() (xPos int, lineNumber int) {
	if cursor.x == 0 {
		return cursor.Position()
	}
	return cursor.x - 1, cursor.line + 1
}

// ForwardCharacter returns the cursors position one character forward.
func (cursor *Cursor) ForwardCharacter() (xPos int, lineNumber int) {
	line, _ := cursor.buffer.GetLine(cursor.line + 1)
	if len(line)-1 < cursor.x {
		return cursor.Position()
	}
	return cursor.x + 1, cursor.line + 1
}

// BeginningOfLine returns the cursors position at the beginning of the line
func (cursor *Cursor) BeginningOfLine() (xPos int, lineNumber int) {
	return 0, cursor.line + 1
}

// EndOfLine returns the cursors position at the end of the line
func (cursor *Cursor) EndOfLine() (xPos int, lineNumber int) {
	line, _ := cursor.buffer.GetLine(cursor.line + 1)
	return len(line) - 1, cursor.line + 1
}

// Pane represents a 'Window' in the editor. It has a Buffer.
type Pane struct {
	buffer  *buffer.Buffer
	cursors map[*buffer.Buffer]*Cursor
	topLine int
}

// NewPane constructs and initilizes a NewPane
func NewPane() Pane {
	return Pane{
		buffer:  nil,
		cursors: make(map[*buffer.Buffer]*Cursor),
		topLine: 1,
	}
}

// Cursor returns the Cursor into the Panes current buffer.
func (pane *Pane) Cursor() *Cursor {
	return pane.cursors[pane.Buffer()]
}

// Buffer returns the Buffer of the Pane
func (pane *Pane) Buffer() *buffer.Buffer {
	return pane.buffer
}

// SetBuffer binds a Buffer to the Pane and creates a Cursor if needed.
func (pane *Pane) SetBuffer(buffer *buffer.Buffer) {
	pane.buffer = buffer
	if pane.Cursor() == nil {
		pane.cursors[pane.buffer] = &Cursor{buffer: pane.buffer}
	}
}

// TopLine returns the line number of the first line visible at the top of the Pane.
func (pane *Pane) TopLine() int {
	return pane.topLine
}

// SetTopLine sets the line number of the first line visiable at the top of the Pane.
func (pane *Pane) SetTopLine(topLine int) {
	pane.topLine = topLine
}

// Editor is the core of Jkl. Maintains buffers, panes and manipluates them.
type Editor struct {
	fa          fileaccessor.FileAccessor
	currentPane *Pane
	buffers     []*buffer.Buffer
	panes       []*Pane
	settings    Settings
}

// New constructs a new editor.
func New(fileaccessor fileaccessor.FileAccessor) Editor {
	pane := NewPane()
	return Editor{
		fa:          fileaccessor,
		currentPane: &pane,
		settings:    DefaultSettings(),
	}
}

// Settings returns the settings.
func (editor *Editor) Settings() *Settings {
	return &editor.settings
}

// OpenFiles opens a list of files into buffers and sets the current buffer to the first of the new buffers.
func (editor *Editor) OpenFiles(filenames []string) {
	for i, filename := range filenames {
		newBuffer := editor.openFile(filename)
		buffer := editor.AddBuffer(&newBuffer)

		if i == 0 {
			editor.CurrentPane().SetBuffer(buffer)
		}
	}
}

// openFile reads a file, loads it into a new buffer and adds it to the list of buffers.
func (editor *Editor) openFile(filename string) buffer.Buffer {
	buffer := buffer.New()

	buffer.SetFilename(filename)

	if data, err := editor.fa.ReadFile(filename); err == nil {
		buffer.SetData(data)
	} else {
		buffer.SetData([]byte{})
	}

	return buffer
}

// OpenFile opens a file and sets it to the current buffer.
func (editor *Editor) OpenFile(filename string) {
	editor.OpenFiles([]string{filename})
}

// AddBuffer adds a buffer to the list of buffers.
func (editor *Editor) AddBuffer(buffer *buffer.Buffer) *buffer.Buffer {
	editor.buffers = append(editor.buffers, buffer)
	return editor.LastBuffer()
}

// LastBuffer returns a pointer to the last buffer in the list of buffers.
func (editor *Editor) LastBuffer() *buffer.Buffer {
	return editor.buffers[len(editor.buffers)-1]
}

// CurrentPane returns the current pane.
func (editor *Editor) CurrentPane() *Pane {
	return editor.currentPane
}

// SetCurrentPane sets the currently visible pane.
func (editor *Editor) SetCurrentPane(pane *Pane) {
	editor.currentPane = pane
}

// Buffers returns a slice containing the buffers.
func (editor *Editor) Buffers() []*buffer.Buffer {
	return editor.buffers
}

// Panes returns a slice containing the panes.
func (editor *Editor) Panes() []*Pane {
	return editor.panes
}
