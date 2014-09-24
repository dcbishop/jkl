package app

import (
	"testing"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/gim/cli"
	. "github.com/smartystreets/goconvey/convey"
)

var fakeFileName = "fakefile.txt"
var fakeFileContents = []byte(`Hello, this is a test`)

var fakeFileSystem = map[string][]byte{
	"fakefile.txt": fakeFileContents,
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

func TestLoadOptions(t *testing.T) {
	app := New(fakeFileAccessor)
	Convey("with default Options doesn't change App", t, func() {
		options := cli.Options{}

		app.LoadOptions(options)

		So(app, ShouldResemble, fakeApp())
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
