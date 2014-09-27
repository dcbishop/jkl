package runegrid

import "github.com/dcbishop/jkl/buffer"

// RuneGrid contains the rendered text UI
type RuneGrid struct {
	width  uint
	height uint
	cells  [][]rune
}

// New constructs a RuneGrid with the given width and height
func New(width, height uint) RuneGrid {
	grid := RuneGrid{
		width:  width,
		height: height,
		cells:  make([][]rune, height),
	}

	for i := range grid.cells {
		grid.cells[i] = make([]rune, width)
	}

	return grid
}

// RenderBuffer blits the buffer onto the grid.
// wrap sets line wrapping on
// linebrake sets soft wrapping
func (grid *RuneGrid) RenderBuffer(
	x, y, x2, y2 uint,
	buffer buffer.Buffer,
	wrap,
	linebrake,
	breakindent bool,
	showbreak string,
	// [TODO]: Should line numbering be done here?
	// Otherwise how do we communicate the line numbers out. - 2014-09-24 11:50pm
) {
	xPos := x
	yPos := y

	for _, r := range buffer.Data() {
		if r == '\n' {
			yPos++
			xPos = x
			continue
		}
		if xPos <= x2 && yPos <= y2 {
			grid.SetCell(xPos, yPos, rune(r))
		}
		xPos++
	}
}

// SetCell sets a cell in the RuneGrid to the given rune
func (grid *RuneGrid) SetCell(x, y uint, r rune) {
	if !grid.IsCellValid(x, y) {
		return
	}

	grid.cells[y][x] = r
}

// IsCellValid returns true if the cell coordinates are valid
func (grid *RuneGrid) IsCellValid(x, y uint) bool {
	if x >= grid.width || y >= grid.height {
		return false
	}
	return true
}

// Width gets the width of the grid.
func (grid *RuneGrid) Width() uint {
	return grid.width
}

// Height gets the height of the grid.
func (grid *RuneGrid) Height() uint {
	return grid.height
}

// Cells gets the cells of the grid.
func (grid *RuneGrid) Cells() [][]rune {
	return grid.cells
}
