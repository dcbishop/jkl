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
	GetLines(first, last int) ([]string, error)
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
	lines, _ := buffer.GetLines(lineNum, lineNum)
	if len(lines) == 1 {
		return lines[0], nil
	}
	return "", errors.New("Not found")
}

// GetLines returns the requested range of lines.
func (buffer *Bytes) GetLines(first, last int) ([]string, error) {
	if first > last {
		return []string{}, errors.New("Invalid range, first > last.")
	}

	lineNum := 1
	pos := 0
	grabbing := false
	lines := []string{}

	for {
		if lineNum >= first {
			grabbing = true
		}

		if grabbing {
			line, err := untillNewLineOrEnd(buffer.data[pos:])
			if err == nil {
				lines = append(lines, line)
			}
		}

		nextLine := bytesUntillNextNewline(buffer.data[pos:]) + 1
		if nextLine == 0 {
			lineNum = last
		}

		if lineNum == last {
			return lines, nil
		}

		pos += nextLine
		lineNum++
	}
}

func untillNewLineOrEnd(data []byte) (string, error) {
	endOfLine := bytesUntillNextNewline(data)
	if endOfLine == -1 {
		if len(data) == 0 {
			return "", errors.New("Not found")
		}
		return string(data), nil
	}
	return string(data[:endOfLine]), nil
}

func bytesUntillNextNewline(data []byte) int {
	return bytes.IndexByte(data, '\n')
}
