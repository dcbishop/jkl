package main

import (
	"strings"
	"testing"

	"github.com/dcbishop/jkl/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

const UnicodeBox = `
	╔═╗
	║ ║
	╚═╝
`

func TestStringToRuneGrid(t *testing.T) {
	Convey("Basic 3x3 unicode box", t, func() {
		expected := [][]rune{
			{'╔', '═', '╗'},
			{'║', 0, '║'},
			{'╚', '═', '╝'},
		}

		runegrid := StringToRuneGrid(UnicodeBox, ' ')
		width, height := runegrid.Size()
		So(width, ShouldEqual, 3)
		So(height, ShouldEqual, 3)
		So(runegrid.cells, ShouldResemble, expected)
	})
}

const Empty3x3 = `
...
...
...
`

func TestRuneGrid(t *testing.T) {
	Convey("New 3x3 grid should be filled with NULL bytes", t, func() {
		grid := NewRuneGrid(3, 3)
		expected := StringToRuneGrid(Empty3x3, '.')

		So(grid, ShouldResemble, expected)

		Convey("SetCell should set the correct cells and ignore invalid ones", func() {
			Compass3x3 := `
			.N.
			WCE
			.S.
			`

			expected := StringToRuneGrid(Compass3x3, '.')
			grid.SetCell(1, 0, 'N')
			grid.SetCell(1, 2, 'S')
			grid.SetCell(2, 1, 'E')
			grid.SetCell(0, 1, 'W')
			grid.SetCell(1, 1, 'C')
			grid.SetCell(100, 100, '!')

			So(grid, ShouldResemble, expected)
		})

		Convey("Basic DrawBox()", func() {
			HashBox := `
			###
			#.#
			###
			`

			expected := StringToRuneGrid(HashBox, '.')

			grid.DrawBox(0, 0, 2, 2, '#')

			So(grid, ShouldResemble, expected)
		})
	})
}

func TestRenderEditor(t *testing.T) {
	Convey("Basic ", t, func() {

		fs := testhelpers.GetTestFs()
		editor := NewEditor(fs)
		editor.OpenFile("file.txt")
		grid := NewRuneGrid(3, 3)

		expected := StringToRuneGrid(Empty3x3, '.')

		So(grid, ShouldResemble, expected)
		Convey("Render 3x3 with default settings", func() {
			UnicodeBox := `
			╔═╗
			║!║
			╚═╝
			`
			expected := StringToRuneGrid(UnicodeBox, 0)
			grid.RenderEditor(&editor)

			So(grid, ShouldResemble, expected)
		})

		Convey("Render 3x3 without borders", func() {
			singleBang := `
			!..
			...
			...
			`
			expected := StringToRuneGrid(singleBang, '.')

			editor.Settings().Borders = false
			grid.RenderEditor(&editor)

			So(grid, ShouldResemble, expected)
		})
	})
}

const OneToNine = `
123
456
789`

var renderTest = []byte(strings.Trim(OneToNine, "\n"))

func TestRenderBuffer(t *testing.T) {
	Convey("Basic Buffer", t, func() {
		buffer := NewBuffer()
		buffer.SetData(renderTest)
		So(CompareBufferBytes(&buffer, renderTest), ShouldBeTrue)

		Convey("3x3 RuneGrid", func() {
			grid := NewRuneGrid(3, 3)
			width, height := grid.Size()
			So(width, ShouldEqual, 3)
			So(height, ShouldEqual, 3)
			Convey("Basic Render", func() {
				expected := StringToRuneGrid(OneToNine, '.')

				settings := DefaultSettings()
				grid.RenderBuffer(&settings, 0, 0, 3, 3, &buffer, 0)

				So(grid, ShouldResemble, expected)
			})

			Convey("offset buffer", func() {
				offsetBuffer := `
				...
				.12
				.45
				`
				expected := StringToRuneGrid(offsetBuffer, '.')

				settings := DefaultSettings()
				grid.RenderBuffer(&settings, 1, 1, 2, 2, &buffer, 1)

				So(grid, ShouldResemble, expected)
			})

			Convey("partially visible buffer", func() {
				partiallyVisible := `
				...
				.1.
				...
				`
				expected := StringToRuneGrid(partiallyVisible, '.')

				settings := DefaultSettings()
				grid.RenderBuffer(&settings, 1, 1, 1, 1, &buffer, 1)

				So(grid, ShouldResemble, expected)
			})

			Convey("larger partially visible buffer", func() {
				partiallyVisible := `
				12.
				45.
				...
				`
				expected := StringToRuneGrid(partiallyVisible, '.')

				settings := DefaultSettings()
				grid.RenderBuffer(&settings, 0, 0, 1, 1, &buffer, 1)

				So(grid, ShouldResemble, expected)
			})
		})
	})

	Convey("Tabs should render as 4 spaces by default", t, func() {
		buffer := NewBuffer()
		buffer.SetData([]byte("{\n\tint a;\n}"))
		tabTest := `
{..........
    int a;.
}..........
...........
`
		grid := NewRuneGrid(11, 4)
		expected := StringToRuneGrid(tabTest, '.')
		settings := DefaultSettings()
		grid.RenderBuffer(&settings, 0, 0, 11, 4, &buffer, 1)
		So(grid, ShouldResemble, expected)
	})
}
