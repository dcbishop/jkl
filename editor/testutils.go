package editor

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

// CompareReaders returns true if they contain the same data.
func CompareReaders(r1 io.Reader, r2 io.Reader) bool {
	b1, err := ioutil.ReadAll(r1)
	if err != nil {
		panic("Error reading reader")
	}

	b2, err := ioutil.ReadAll(r2)
	if err != nil {
		panic("Error reading reader")
	}

	return bytes.Compare(b1, b2) == 0
}

// CompareBufferString returns true if the contents of buffer matches those in the string.
func CompareBufferString(buffer *Buffer, s string) bool {
	r := strings.NewReader(s)
	return CompareReaders(r, buffer.NewReader())
}

// CompareBufferBytes returns true if the contents of buffer matches those in the string.
func CompareBufferBytes(buffer *Buffer, b []byte) bool {
	s := string(b)
	return CompareBufferString(buffer, s)
}
