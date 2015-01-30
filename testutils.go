package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/spf13/afero"
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

// A Unicode box
const UnicodeBox = `
	╔═╗
	║ ║
	╚═╝
`

// A 3x3 matrix of .'s
const Empty3x3 = `
...
...
...
`

// StringToRunes converts ASCII art to a RuneGrid for testing.
func StringToRunes(s string, replaceWithNul rune) [][]rune {
	s = strings.Trim(s, "\n\t ")
	s = strings.Replace(s, string(replaceWithNul), string(0), 9999)
	s = strings.Replace(s, "\t", "", 9999)

	// Get width based on number of characters in the first line
	// (strings.Index doesn't seem to work with the unicode box example)
	runes := []rune(s)
	i := 0
	for runes[i] != '\n' {
		i++
	}
	width := i

	// Get height based on line numbers
	height := strings.Count(s, "\n") + 1

	// Strip newlines
	runes = []rune(strings.Replace(s, "\n", "", 9999))

	i = 0

	grid := make([][]rune, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			r := runes[i]
			i++
			grid[y][x] = r
		}
	}
	return grid
}

// StringToRuneGrid converts ASCII art to a RuneGrid for testing.
func StringToRuneGrid(s string, replaceWithNul rune) RuneGrid {
	cells := StringToRunes(s, replaceWithNul)
	height := len(cells)
	if height == 1 {
		return NewRuneGrid(0, 0)
	}
	width := len(cells[0])
	grid := NewRuneGrid(width, height)
	grid.cells = cells
	return grid
}

var fakeFileName = "fakefile.txt"
var fakeFileContents = []byte(`Hello, this is a test`)
var fakeFileContents2 = []byte(`This is the 2nd file!`)

var fakeFileSystem = map[string][]byte{
	"file.txt":      []byte(`!`),
	"fakefile.txt":  fakeFileContents,
	"fakefile2.txt": fakeFileContents2,
}

// GetTestFs returns an inmemory filesystem with test data.
func GetTestFs() afero.Fs {
	return GetCustomTestFs(fakeFileSystem)
}

// GetCustomTestFs returns an inmemory filesystem with the given test files.
func GetCustomTestFs(filemap map[string][]byte) afero.Fs {
	fs := afero.MemMapFs{}

	for filename, data := range filemap {
		file, _ := fs.Create(filename)
		defer file.Close()
		file.Write(data)
	}

	return &fs
}
