package editor

import (
	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/buffer"
)

// Editor contains buffers and performs actions on them.
type Editor interface {
	OpenFiles(filenames []string)
	OpenFile(filename string)
	AddBuffer(buffer buffer.Buffer) buffer.Buffer
	LastBuffer() buffer.Buffer
	CurrentBuffer() buffer.Buffer
	SetCurrentBuffer(buffer buffer.Buffer)
	Buffers() []buffer.Buffer
}

// Jkl is the standard implementation of Editor
type Jkl struct {
	buffers       []buffer.Buffer
	currentBuffer buffer.Buffer
	fa            fileaccessor.FileAccessor
}

// New constructs a new editor
func New(fileaccessor fileaccessor.FileAccessor) Jkl {
	return Jkl{fa: fileaccessor}
}

// OpenFiles opens a list of files into buffers and sets the current buffer to the first of the new buffers.
func (editor *Jkl) OpenFiles(filenames []string) {
	for i, filename := range filenames {
		buffer := editor.openFile(filename)
		buffer = editor.AddBuffer(buffer)

		if i == 0 {
			editor.SetCurrentBuffer(buffer)
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

// CurrentBuffer returns the current buffer.
func (editor *Jkl) CurrentBuffer() buffer.Buffer {
	return editor.currentBuffer
}

// SetCurrentBuffer sets the currently visible buffer.
func (editor *Jkl) SetCurrentBuffer(buffer buffer.Buffer) {
	editor.currentBuffer = buffer
}

// Buffers returns a slice containing the buffers.
func (editor *Jkl) Buffers() []buffer.Buffer {
	return editor.buffers
}
