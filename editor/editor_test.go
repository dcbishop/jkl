package editor

import (
	"testing"

	"github.com/dcbishop/jkl/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

var fakeFileName = "fakefile.txt"
var fakeFileContents = []byte(`Hello, this is a test`)
var fakeFileContents2 = []byte(`This is the 2nd file!`)

var fakeFileSystem = map[string][]byte{
	"fakefile.txt":  fakeFileContents,
	"fakefile2.txt": fakeFileContents,
}

func TestOpenFile(t *testing.T) {
	Convey("editor.OpenFile", t, func() {
		fs := testhelpers.GetCustomTestFs(fakeFileSystem)
		editor := New(fs)
		Convey("with a valid filename, loads file into buffer", func() {
			So(len(editor.buffers), ShouldEqual, 0)

			editor.OpenFile(fakeFileName)

			So(len(editor.Buffers()), ShouldEqual, 1)
			So(CompareBufferBytes(editor.buffers[0], fakeFileContents), ShouldBeTrue)
			So(editor.buffers[0].Filename(), ShouldResemble, fakeFileName)
		})

		Convey("with nonexistant file, opens a blank buffer with that filename", func() {
			So(len(editor.Buffers()), ShouldEqual, 0)
			editor.OpenFile("file2.txt")
			So(len(editor.Buffers()), ShouldEqual, 1)
			So(CompareBufferBytes(editor.buffers[0], []byte{}), ShouldBeTrue)
			So(editor.buffers[0].Filename(), ShouldResemble, "file2.txt")
		})
	})
}

var data = []byte(`Hello, World!
Line 2
Line 3`)

func TestCursor(t *testing.T) {
	Convey("New cursor", t, func() {
		buffer := NewBuffer()
		buffer.SetData(data)
		cursor := Cursor{buffer: &buffer}

		x, line := cursor.Position()

		So(x, ShouldEqual, 0)
		So(line, ShouldEqual, 1)
		So(cursor.buffer, ShouldEqual, &buffer)

		Convey("DownLine", func() {
			x, line := cursor.DownLine()
			So(x, ShouldEqual, 0)
			So(line, ShouldEqual, 2)

			Convey("Move DownLine", func() {
				cursor.Move(cursor.DownLine())
				x, line := cursor.Position()
				So(x, ShouldEqual, 0)
				So(line, ShouldEqual, 2)

				Convey("UpLine", func() {
					x, line := cursor.UpLine()
					So(x, ShouldEqual, 0)
					So(line, ShouldEqual, 1)
				})
			})

		})

		Convey("UpLine when at top", func() {
			x, line := cursor.UpLine()
			So(x, ShouldEqual, 0)
			So(line, ShouldEqual, 1)
		})

		Convey("ForwardCharacter should move forward one character", func() {
			x, _ := cursor.ForwardCharacter()
			So(x, ShouldEqual, 1)

			Convey("BackCharacter should move backward one character", func() {
				cursor.Move(cursor.ForwardCharacter())
				So(x, ShouldEqual, 1)
				x, _ := cursor.BackCharacter()
				So(x, ShouldEqual, 0)
			})
		})

		Convey("BackCharacter when at first character shouldn't move", func() {
			x, _ := cursor.BackCharacter()
			So(x, ShouldEqual, 0)
		})

		Convey("Move to EndOfLine", func() {
			cursor.Move(cursor.EndOfLine())
			x, line := cursor.Position()
			So(x, ShouldEqual, 12)
			So(line, ShouldEqual, 1)

			Convey("BeginningOfLine", func() {
				x, line := cursor.BeginningOfLine()
				So(x, ShouldEqual, 0)
				So(line, ShouldEqual, 1)
			})
		})
	})
}
