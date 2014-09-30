package testhelpers

import "strings"

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
