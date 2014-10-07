package buffer

import (
	"bytes"
	"errors"
)

// Buffer contains the text to edit.
type Buffer interface {
	Filename() string
	SetFilename(filename string)
	Data() []byte
	SetData(data []byte)
	GetLine(lineNum int) (string, error)
}

// Bytes is a simple implenetation of the buffer that stores data in a []byte.
type Bytes struct {
	filename string
	data     []byte
}

// New constructs a new ByteBuffer object containing data.
func New() Bytes {
	return Bytes{}
}

// Filename returns the buffers filename as a byte slice.
func (buffer *Bytes) Filename() string {
	return buffer.filename
}

// SetFilename ses the filename of that will be used when the buffer is saved.
func (buffer *Bytes) SetFilename(filename string) {
	buffer.filename = filename
}

// Data returns the buffers data as a byte slice.
func (buffer *Bytes) Data() []byte {
	return buffer.data
}

// SetData sets the data of the buffer.
func (buffer *Bytes) SetData(data []byte) {
	buffer.data = data
}

// GetLine returns the requested line as a string.
func (buffer *Bytes) GetLine(lineNum int) (string, error) {
	if lineNum < 0 {
		return "", errors.New("Not found")
	}

	line := 1
	pos := 0
	for {
		if line == lineNum {
			return untillNewLineOrEnd(buffer.data[pos:]), nil
		}

		nextLine := bytesUntillNextNewline(buffer.data[pos:]) + 1
		if nextLine == 0 {
			return "", errors.New("Not found")
		}

		pos += nextLine
		line++
	}
}

func untillNewLineOrEnd(data []byte) string {
	endOfLine := bytesUntillNextNewline(data)
	if endOfLine == -1 {
		return string(data)
	}
	return string(data[:endOfLine])
}

func bytesUntillNextNewline(data []byte) int {
	return bytes.IndexByte(data, '\n')
}
