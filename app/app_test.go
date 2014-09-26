package app

import (
	"testing"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/gim/cli"
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
	Convey("app.Run() should be terminated by an App.Stop() in a reasonable time", t, func() {
		app := fakeApp()
		time.AfterFunc(time.Second/5, func() {
			panic("Failed to stop app.")
		})
		go app.Stop()
		app.Run()
	})
	Convey("app.Stop() without having called app.Run() shouldn't explode", t, func() {
		app := fakeApp()
		app.Stop()
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
		options, err := cli.ParseArgs([]string{"gim", "fakefile.txt", "fakefile2.txt"})
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
