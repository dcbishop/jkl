package main

import (
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuffer(t *testing.T) {
	Convey("NewBuffer() buffer", t, func() {
		buffer := NewBuffer()
		r := buffer.NewReader()
		bytes, _ := ioutil.ReadAll(r)
		So(len(bytes), ShouldEqual, 0)
		So(buffer.Filename(), ShouldEqual, "")

		Convey("SetFilename()", func() {
			buffer.SetFilename("file.txt")
			So(buffer.Filename(), ShouldEqual, "file.txt")
		})

		Convey("SetDataString()", func() {
			data := "Hello, World!"
			buffer.SetDataString(data)
			So(CompareBufferString(&buffer, data), ShouldBeTrue)
		})
	})
}

var testData = `This is line 1.
This is line 2.
This is line 3.`

func TestGetLine(t *testing.T) {
	Convey("Test multiline Buffer", t, func() {
		buffer := NewBuffer()
		buffer.SetDataString(testData)

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

	Convey("Test line feed terminated Buffer", t, func() {
		// "Test" with single line feed
		data := []byte{84, 101, 115, 116, 10}

		Convey("get first line", func() {
			buffer := NewBuffer()
			buffer.SetData(data)
			line, err := buffer.GetLine(1)
			So(line, ShouldResemble, "Test")
			So(err, ShouldBeNil)
		})

		Convey("get missing 2nd line", func() {
			buffer := NewBuffer()
			buffer.SetData(data)
			line, err := buffer.GetLine(2)
			So(line, ShouldResemble, "")
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Test double line feed terminated Buffer", t, func() {
		// "Test" with double line feed
		data := []byte{84, 101, 115, 116, 10, 10}
		buffer := NewBuffer()
		buffer.SetData(data)
		line, err := buffer.GetLine(2)
		So(line, ShouldResemble, "")
		So(err, ShouldBeNil)
	})
}

var testData2 = `1
2
3`

func TestGetLines(t *testing.T) {
	Convey("Test multiline Buffer", t, func() {
		buffer := NewBuffer()
		buffer.SetDataString(testData2)

		Convey("get all lines", func() {
			line, err := buffer.GetLines(1, 3)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, []string{"1", "2", "3"})
		})

		Convey("get first 2 lines", func() {
			line, err := buffer.GetLines(1, 2)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, []string{"1", "2"})
		})

		Convey("get last 2 lines", func() {
			line, err := buffer.GetLines(2, 3)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, []string{"2", "3"})
		})

		Convey("get single line", func() {
			line, err := buffer.GetLines(2, 2)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, []string{"2"})
		})

		Convey("get reversed range should return an error", func() {
			line, err := buffer.GetLines(3, 1)
			So(err, ShouldNotBeNil)
			So(line, ShouldResemble, []string{})
		})

		Convey("get more lines than exist should get lines that do exist", func() {
			line, err := buffer.GetLines(1, 4)
			So(err, ShouldBeNil)
			So(line, ShouldResemble, []string{"1", "2", "3"})
		})

		Convey("get negative line should get lines that do exist", func() {
			line, err := buffer.GetLines(-1, 3)
			So(line, ShouldResemble, []string{"1", "2", "3"})
			So(err, ShouldBeNil)
		})
	})
}
