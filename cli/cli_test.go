package cli

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseArgs(t *testing.T) {
	Convey("ParseArgs", t, func() {
		Convey("with an empty slice", func() {
			result, err := ParseArgs([]string{})
			So(err, ShouldBeNil)
			So(result, ShouldResemble, Options{})
		})
		Convey("with a single filename", func() {
			result, err := ParseArgs([]string{"gim", "file.txt"})
			So(err, ShouldBeNil)
			So(result, ShouldResemble, Options{FilesToOpen: []string{"file.txt"}})
		})
		Convey("with multiple files", func() {
			result, err := ParseArgs([]string{"gim", "file1.txt", "file2.txt", "file3.txt"})
			So(err, ShouldBeNil)
			So(result, ShouldResemble, Options{FilesToOpen: []string{"file1.txt", "file2.txt", "file3.txt"}})
		})
		Convey("with --help", func() {
			result, err := ParseArgs([]string{"gim", "--help"})
			So(err, ShouldBeNil)
			So(result, ShouldResemble, Options{Help: true})
		})
	})
}
