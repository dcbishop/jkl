package editor

import (
	"testing"

	"github.com/dcbishop/fileaccessor"
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

func TestOpenFile(t *testing.T) {
	Convey("editor.OpenFile", t, func() {
		editor := New(fakeFileAccessor)
		Convey("with a valid filename, loads file into buffer", func() {
			So(len(editor.buffers), ShouldEqual, 0)

			editor.OpenFile(fakeFileName)

			So(len(editor.Buffers()), ShouldEqual, 1)
			So(editor.buffers[0].Data(), ShouldResemble, fakeFileContents)
			So(editor.buffers[0].Filename(), ShouldResemble, fakeFileName)
		})

		Convey("with nonexistant file, opens a blank buffer with that filename", func() {
			So(len(editor.Buffers()), ShouldEqual, 0)
			editor.OpenFile("file2.txt")
			So(len(editor.Buffers()), ShouldEqual, 1)
			So(editor.buffers[0].Data(), ShouldResemble, []byte{})
			So(editor.buffers[0].Filename(), ShouldResemble, "file2.txt")
		})
	})
}
