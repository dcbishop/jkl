package editor

import (
	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/buffer"
)

// Settings stores settings for the editor
type Settings struct {
	Borders     bool
	OuterBorder bool
}

// DefaultSettings constructs a default settings.
func DefaultSettings() Settings {
	return Settings{
		Borders:     true,
		OuterBorder: true,
	}
}

// Cursor stores a position in a buffer and handles movement.
type Cursor struct {
	x      int
	line   int
	buffer buffer.Buffer
}

// Position reutrns the cursors current position.
func (cursor *Cursor) Position() (xPos int, lineNumber int) {
	return cursor.x, cursor.line
}

// Move repositions the cursor at the given coordinates.
func (cursor *Cursor) Move(xPos int, lineNumber int) {
	cursor.x = xPos
	cursor.line = lineNumber
}

// DownLine returns the cursors position one line down.
func (cursor *Cursor) DownLine() (xPos int, lineNumber int) {
	return cursor.x, cursor.line + 1
}

// UpLine returns the cursors position one line up.
func (cursor *Cursor) UpLine() (xPos int, lineNumber int) {
	if cursor.line == 0 {
		return cursor.Position()
	}
	return cursor.x, cursor.line - 1
}

// BackCharacter returns the cursors position one character back.
func (cursor *Cursor) BackCharacter() (xPos int, lineNumber int) {
	if cursor.x == 0 {
		return cursor.Position()
	}
	return cursor.x - 1, cursor.line
}

// ForwardCharacter returns the cursors position one character forward.
func (cursor *Cursor) ForwardCharacter() (xPos int, lineNumber int) {
	return cursor.x + 1, cursor.line
}

// Pane represents a 'Window' in the editor. It has a Buffer.
type Pane struct {
	buffer  buffer.Buffer
	cursors map[buffer.Buffer]*Cursor
}

// NewPane constructs and initilizes a NewPane
func NewPane() Pane {
	return Pane{
		buffer:  nil,
		cursors: make(map[buffer.Buffer]*Cursor),
	}
}

// Cursor returns the Cursor into the Panes current buffer.
func (pane *Pane) Cursor() *Cursor {
	return pane.cursors[pane.Buffer()]
}

// Buffer returns the Buffer of the Pane
func (pane *Pane) Buffer() buffer.Buffer {
	return pane.buffer
}

// SetBuffer binds a Buffer to the Pane and creates a Cursor if needed.
func (pane *Pane) SetBuffer(buffer buffer.Buffer) {
	pane.buffer = buffer
	if pane.Cursor() == nil {
		pane.cursors[pane.buffer] = &Cursor{buffer: pane.buffer}
	}
}

// Editor contains buffers and performs actions on them.
type Editor interface {
	OpenFiles(filenames []string)
	OpenFile(filename string)
	AddBuffer(buffer buffer.Buffer) buffer.Buffer
	CurrentPane() *Pane
	SetCurrentPane(pane *Pane)
	Buffers() []buffer.Buffer
	LastBuffer() buffer.Buffer
	Panes() []*Pane
	Settings() *Settings
}

// Jkl is the standard implementation of Editor
type Jkl struct {
	fa          fileaccessor.FileAccessor
	currentPane *Pane
	buffers     []buffer.Buffer
	panes       []*Pane
	settings    Settings
}

// New constructs a new editor
func New(fileaccessor fileaccessor.FileAccessor) Jkl {
	pane := NewPane()
	return Jkl{
		fa:          fileaccessor,
		currentPane: &pane,
		settings:    DefaultSettings(),
	}
}

// Settings returns the settings
func (editor *Jkl) Settings() *Settings {
	return &editor.settings
}

// OpenFiles opens a list of files into buffers and sets the current buffer to the first of the new buffers.
func (editor *Jkl) OpenFiles(filenames []string) {
	for i, filename := range filenames {
		buffer := editor.openFile(filename)
		buffer = editor.AddBuffer(buffer)

		if i == 0 {
			editor.CurrentPane().SetBuffer(buffer)
		}
	}
}

// openFile reads a file, loads it into a new buffer and adds it to the list of buffers
func (editor *Jkl) openFile(filename string) buffer.Buffer {
	buffer := buffer.New()

	buffer.SetFilename(filename)

	if data, err := editor.fa.ReadFile(filename); err == nil {
		buffer.SetData(data)
	} else {
		buffer.SetData([]byte{})
	}

	return &buffer
}

// OpenFile opens a file and sets it to the current buffer.
func (editor *Jkl) OpenFile(filename string) {
	editor.OpenFiles([]string{filename})
}

// AddBuffer adds a buffer to the list of buffers
func (editor *Jkl) AddBuffer(buffer buffer.Buffer) buffer.Buffer {
	editor.buffers = append(editor.buffers, buffer)
	return editor.LastBuffer()
}

// LastBuffer returns a pointer to the last buffer in the list of buffers
func (editor *Jkl) LastBuffer() buffer.Buffer {
	return editor.buffers[len(editor.buffers)-1]
}

// CurrentPane returns the current pane.
func (editor *Jkl) CurrentPane() *Pane {
	return editor.currentPane
}

// SetCurrentPane sets the currently visible pane.
func (editor *Jkl) SetCurrentPane(pane *Pane) {
	editor.currentPane = pane
}

// Buffers returns a slice containing the buffers.
func (editor *Jkl) Buffers() []buffer.Buffer {
	return editor.buffers
}

// Panes returns a slice containing the panes.
func (editor *Jkl) Panes() []*Pane {
	return editor.panes
}
