package app

import (
	"testing"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/cli"
	"github.com/dcbishop/jkl/service"
	. "github.com/smartystreets/goconvey/convey"
)

var fakeFileName = "fakefile.txt"
var fakeFileContents = []byte(`Hello, this is a test`)
var fakeFileContents2 = []byte(`This is the 2nd file!`)

var fakeFileSystem = map[string][]byte{
	"fakefile.txt":  fakeFileContents,
	"fakefile2.txt": fakeFileContents,
}

var fakeFileAccessor = fileaccessor.Virtual{fakeFileSystem}

func fakeApp() App {
	return New(fakeFileAccessor)
}

func TestNew(t *testing.T) {
	Convey("app.New() returns basic App", t, func() {
		app := fakeApp()
		So(app, ShouldResemble, fakeApp())
	})
}

func TestRunStop(t *testing.T) {
	Convey("given a new app", t, func() {
		app := fakeApp()
		So(app.Running(), ShouldBeFalse)

		Convey("app.Run() should start the app in a reasonable time", func() {
			go app.Run()
			So(service.WaitUntilRunning(&app, time.Second), ShouldBeNil)
			So(app.Running(), ShouldBeTrue)

			Convey("app.Stop() should terminate in a reasonable time", func() {
				go app.Stop()
				So(service.WaitUntilStopped(&app, time.Second), ShouldBeNil)
				So(app.Running(), ShouldBeFalse)
			})

			Convey("app.Run() on a running app should panic", func() {
				So(app.Run, ShouldPanic)
			})
		})

		Convey("app.Stop() without having called app.Run() shouldn't explode", func() {
			app.Stop()
			So(app.Running(), ShouldBeFalse)
		})
	})
}

func TestLoadOptions(t *testing.T) {
	app := New(fakeFileAccessor)
	Convey("with default Options doesn't change App", t, func() {
		options := cli.Options{}

		app.LoadOptions(options)

		So(app, ShouldResemble, fakeApp())
	})
	Convey("with 2 filenames given", t, func() {
		options, err := cli.ParseArgs([]string{"jkl", "fakefile.txt", "fakefile2.txt"})
		So(err, ShouldBeNil)

		app.LoadOptions(options)
		So(len(app.buffers), ShouldEqual, 2)
	})
}

func TestOpenFile(t *testing.T) {
	Convey("app.OpenFile", t, func() {
		app := fakeApp()
		Convey("with a valid filename, loads file into buffer", func() {
			So(len(app.buffers), ShouldEqual, 0)

			app.OpenFile(fakeFileName)

			So(len(app.buffers), ShouldEqual, 1)
			So(app.buffers[0].data, ShouldResemble, fakeFileContents)
			So(app.buffers[0].filename, ShouldResemble, fakeFileName)
		})

		Convey("with nonexistant file, opens a blank buffer with that filename", func() {
			So(len(app.buffers), ShouldEqual, 0)
			app.OpenFile("file2.txt")
			So(len(app.buffers), ShouldEqual, 1)
			So(app.buffers[0].data, ShouldResemble, []byte{})
			So(app.buffers[0].filename, ShouldResemble, "file2.txt")
		})
	})
}

func TestRuneGrid(t *testing.T) {
	Convey("New grid should be filled with NULL bytes", t, func() {
		grid := NewRuneGrid(3, 3)
		expected := [][]rune{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		}
		So(grid.cells, ShouldResemble, expected)
	})
	Convey("SetCell should set the correct cells and ignore invalid ones", t, func() {
		grid := NewRuneGrid(3, 3)
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
}

var renderTest = []byte(`123
456
789`)

func TestRenderBuffer(t *testing.T) {
	Convey("Basic Buffer", t, func() {
		buffer := NewBuffer()
		buffer.data = renderTest
		grid := NewRuneGrid(3, 3)
		expected := [][]rune{
			{'1', '2', '3'},
			{'4', '5', '6'},
			{'7', '8', '9'},
		}

		grid.RenderBuffer(0, 0, 3, 3, &buffer, false, false, false, "")

		So(grid.cells, ShouldResemble, expected)
	})
	Convey("offset buffer", t, func() {
		buffer := NewBuffer()
		buffer.data = renderTest
		grid := NewRuneGrid(3, 3)
		expected := [][]rune{
			{0, 0, 0},
			{0, '1', '2'},
			{0, '4', '5'},
		}

		grid.RenderBuffer(1, 1, 2, 2, &buffer, false, false, false, "")

		So(grid.cells, ShouldResemble, expected)
	})
	Convey("partially visible buffer", t, func() {
		buffer := NewBuffer()
		buffer.data = renderTest
		grid := NewRuneGrid(3, 3)
		expected := [][]rune{
			{0, 0, 0},
			{0, '1', 0},
			{0, 0, 0},
		}

		grid.RenderBuffer(1, 1, 1, 1, &buffer, false, false, false, "")

		So(grid.cells, ShouldResemble, expected)
	})
	Convey("larger partially visible buffer", t, func() {
		buffer := NewBuffer()
		buffer.data = renderTest
		grid := NewRuneGrid(3, 3)
		expected := [][]rune{
			{'1', '2', 0},
			{'4', '5', 0},
			{0, 0, 0},
		}

		grid.RenderBuffer(0, 0, 1, 1, &buffer, false, false, false, "")

		So(grid.cells, ShouldResemble, expected)
	})
}
