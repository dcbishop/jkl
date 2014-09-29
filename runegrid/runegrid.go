package runegrid

import (
	"github.com/dcbishop/jkl/buffer"
	"github.com/dcbishop/jkl/editor"
)

// RuneGrid contains the rendered text UI
type RuneGrid struct {
	width  int
	height int
	cells  [][]rune
}

// New constructs a RuneGrid with the given width and height
func New(width, height int) RuneGrid {
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

// RenderEditor renders the entire editor window to the grid.
func (grid *RuneGrid) RenderEditor(editor editor.Editor) {
	x1 := 0
	y1 := 0
	x2 := grid.width - 1
	y2 := grid.height - 1

	settings := editor.Settings()

	if settings.Borders && settings.OuterBorder {
		grid.DrawBox(x1, y1, x2, y2, '═', '║', '╔', '╗', '╚', '╝')

		x1++
		y1++
		x2--
		y2--
	}

	if editor.CurrentPane() == nil {
		return
	}

	grid.RenderPane(editor.CurrentPane(), x1, y1, x2, y2)
}

// RenderPane render the Pane and it's contents
func (grid *RuneGrid) RenderPane(pane *editor.Pane, x1, y1, x2, y2 int) {
	if pane.Buffer() == nil {
		return
	}
	grid.RenderBuffer(x1, y1, x2, y2, pane.Buffer(), false, false, false, "")
}

// RenderBuffer blits the buffer onto the grid.
// wrap sets line wrapping on
// linebrake sets soft wrapping
func (grid *RuneGrid) RenderBuffer(
	x1, y1, x2, y2 int,
	buffer buffer.Buffer,
	wrap,
	linebrake,
	breakindent bool,
	showbreak string,
	// [TODO]: Should line numbering be done here?
	// Otherwise how do we communicate the line numbers out. - 2014-09-24 11:50pm
) {
	xPos := x1
	yPos := y1

	for _, r := range buffer.Data() {
		if r == '\n' {
			yPos++
			xPos = x1
			continue
		}
		if xPos <= x2 && yPos <= y2 {
			grid.SetCell(xPos, yPos, rune(r))
		}
		xPos++
	}
}

// SetCell sets a cell in the RuneGrid to the given rune
func (grid *RuneGrid) SetCell(x, y int, r rune) {
	if !grid.IsCellValid(x, y) {
		return
	}

	grid.cells[y][x] = r
}

// IsCellValid returns true if the cell coordinates are valid
func (grid *RuneGrid) IsCellValid(x, y int) bool {
	if x >= grid.width || y >= grid.height {
		return false
	}
	return true
}

// Width gets the width of the grid.
func (grid *RuneGrid) Width() int {
	return grid.width
}

// Height gets the height of the grid.
func (grid *RuneGrid) Height() int {
	return grid.height
}

// Cells gets the cells of the grid.
func (grid *RuneGrid) Cells() [][]rune {
	return grid.cells
}

// DrawBox a box with the given runes.
func (grid *RuneGrid) DrawBox(x1, y1, x2, y2 int, r rune, rExtra ...rune) {
	if len(rExtra) != 0 && len(rExtra) != 1 && len(rExtra) != 5 {
		panic("rExtra must be 0, 1 or 5 arguments")
	}

	vertical := r
	horizontal := r
	topLeft := r
	topRight := r
	bottomLeft := r
	bottomRight := r

	var _ = vertical
	var _ = horizontal
	var _ = topLeft
	var _ = topRight
	var _ = bottomLeft
	var _ = bottomRight

	if len(rExtra) > 0 {
		vertical = rExtra[0]
	}

	if len(rExtra) == 5 {
		topLeft = rExtra[1]
		topRight = rExtra[2]
		bottomLeft = rExtra[3]
		bottomRight = rExtra[4]
	}

	grid.DrawHorizontalLine(x1+1, x2-1, y1, horizontal)
	grid.DrawHorizontalLine(x1+1, x2-1, y2, horizontal)
	grid.DrawVerticalLine(x1, y1+1, y2-1, vertical)
	grid.DrawVerticalLine(x2, y1+1, y2-1, vertical)

	grid.SetCell(x1, y1, topLeft)
	grid.SetCell(x2, y1, topRight)
	grid.SetCell(x1, y2, bottomLeft)
	grid.SetCell(x2, y2, bottomRight)

}

// DrawHorizontalLine draws a line with the given rune
func (grid *RuneGrid) DrawHorizontalLine(x1, x2, y int, r rune) {
	for x := x1; x <= x2; x++ {
		grid.SetCell(x, y, r)
	}
}

// DrawVerticalLine draws a vertical line
func (grid *RuneGrid) DrawVerticalLine(x, y1, y2 int, r rune) {
	for y := y1; y <= y2; y++ {
		grid.SetCell(x, y, r)
	}
}
