package buffer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuffer(t *testing.T) {
	Convey("New() buffer", t, func() {
		buffer := New()
		So(len(buffer.Data()), ShouldEqual, 0)
		So(buffer.Filename(), ShouldEqual, "")

		Convey("SetFilename()", func() {
			buffer.SetFilename("file.txt")
			So(buffer.Filename(), ShouldEqual, "file.txt")
		})

		Convey("SetData()", func() {
			data := []byte("Hello, World!")
			buffer.SetData(data)
			So(buffer.Data(), ShouldResemble, data)
		})
	})
}

var testData = []byte(`This is line 1.
This is line 2.
This is line 3.`)

func TestGetLine(t *testing.T) {
	Convey("OnTestBuffer", t, func() {
		buffer := New()
		buffer.SetData(testData)

		Convey("get first line", func() {
			line, err := buffer.GetLine(1)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, "This is line 1.")
		})
		Convey("get second line", func() {
			line, err := buffer.GetLine(2)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, "This is line 2.")
		})
		Convey("get third line", func() {
			line, err := buffer.GetLine(3)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, "This is line 3.")
		})
		Convey("get missing forth line", func() {
			line, err := buffer.GetLine(4)
			So(line, ShouldResemble, "")
			So(err, ShouldNotBeNil)
		})
		Convey("get -1 line", func() {
			line, err := buffer.GetLine(-1)
			So(line, ShouldResemble, "")
			So(err, ShouldNotBeNil)
		})
	})
}
