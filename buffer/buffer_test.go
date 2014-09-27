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
