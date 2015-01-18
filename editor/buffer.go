package editor

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// Buffer is a simple implenetation of the buffer that stores data in a []byte.
type Buffer struct {
	filename string
	data     []byte
}

// NewBuffer constructs a new ByteBuffer object containing data.
func NewBuffer() Buffer {
	return Buffer{}
}

// Filename returns the buffers filename as a byte slice.
func (buffer *Buffer) Filename() string {
	return buffer.filename
}

// SetFilename ses the filename of that will be used when the buffer is saved.
func (buffer *Buffer) SetFilename(filename string) {
	buffer.filename = filename
}

// NewReader returns a new Reader reading from the buffer contents.
func (buffer *Buffer) NewReader() io.Reader {
	b := bytes.NewReader(buffer.data)
	return b
}

// SetData sets the data of the buffer.
func (buffer *Buffer) SetData(data []byte) {
	buffer.ReadData(bytes.NewReader(data))
}

// SetDataString sets the data of the buffer.
func (buffer *Buffer) SetDataString(data string) {
	buffer.ReadData(strings.NewReader(data))
}

// ReadData reads the data from the given stream into the buffer replacing anything already there.
func (buffer *Buffer) ReadData(data io.Reader) {
	b := new(bytes.Buffer)
	b.ReadFrom(data)
	buffer.setData(b.Bytes())
}

func (buffer *Buffer) setData(data []byte) {
	buffer.data = data
}

// GetLine returns the requested line as a string.
func (buffer *Buffer) GetLine(lineNum int) (string, error) {
	lines, _ := buffer.GetLines(lineNum, lineNum)
	if len(lines) == 1 {
		return lines[0], nil
	}
	return "", errors.New("Not found")
}

// GetLines returns the requested range of lines.
func (buffer *Buffer) GetLines(first, last int) ([]string, error) {
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
