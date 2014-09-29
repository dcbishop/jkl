package runegrid

import (
	"testing"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/buffer"
	"github.com/dcbishop/jkl/editor"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRuneGrid(t *testing.T) {
	Convey("New 3x3 grid should be filled with NULL bytes", t, func() {
		grid := New(3, 3)
		expected := [][]rune{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}

		So(grid.cells, ShouldResemble, expected)
		Convey("SetCell should set the correct cells and ignore invalid ones", func() {
			grid := New(3, 3)
			expected := [][]rune{
				{0, 'N', 0},
				{'W', 'C', 'E'},
				{0, 'S', 0},
			}
			grid.SetCell(1, 0, 'N')
			grid.SetCell(1, 2, 'S')
			grid.SetCell(2, 1, 'E')
			grid.SetCell(0, 1, 'W')
			grid.SetCell(1, 1, 'C')
			grid.SetCell(100, 100, '!')

			So(grid.cells, ShouldResemble, expected)
		})
		Convey("Basic DrawBox()", func() {
			grid.DrawBox(0, 0, 2, 2, '#')
			expected := [][]rune{
				{'#', '#', '#'},
				{'#', 0, '#'},
				{'#', '#', '#'},
			}
			So(grid.cells, ShouldResemble, expected)
		})
	})
}

func TestRenderEditor(t *testing.T) {
	Convey("Basic Editor", t, func() {
		fa := fileaccessor.Virtual{}
		editor := editor.New(fa)
		grid := New(3, 3)
		grid.RenderEditor(&editor)
	})
}

var renderTest = []byte(`123
456
789`)

func TestRenderBuffer(t *testing.T) {
	Convey("Basic Buffer", t, func() {
		buffer := buffer.New()
		buffer.SetData(renderTest)
		Convey("3x3 RuneGrid", func() {
			grid := New(3, 3)
			So(grid.Width(), ShouldEqual, 3)
			So(grid.Height(), ShouldEqual, 3)
			Convey("Basic Render", func() {
				expected := [][]rune{
					{'1', '2', '3'},
					{'4', '5', '6'},
					{'7', '8', '9'},
				}

				grid.RenderBuffer(0, 0, 3, 3, &buffer, false, false, false, "")

				So(grid.Cells(), ShouldResemble, expected)
			})
			Convey("offset buffer", func() {
				expected := [][]rune{
					{0, 0, 0},
					{0, '1', '2'},
					{0, '4', '5'},
				}

				grid.RenderBuffer(1, 1, 2, 2, &buffer, false, false, false, "")

				So(grid.Cells(), ShouldResemble, expected)
			})
			Convey("partially visible buffer", func() {
				expected := [][]rune{
					{0, 0, 0},
					{0, '1', 0},
					{0, 0, 0},
				}

				grid.RenderBuffer(1, 1, 1, 1, &buffer, false, false, false, "")

				So(grid.Cells(), ShouldResemble, expected)
			})
			Convey("larger partially visible buffer", func() {
				expected := [][]rune{
					{'1', '2', 0},
					{'4', '5', 0},
					{0, 0, 0},
				}

				grid.RenderBuffer(0, 0, 1, 1, &buffer, false, false, false, "")

				So(grid.Cells(), ShouldResemble, expected)
			})
		})
	})
}
