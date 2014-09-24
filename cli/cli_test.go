package cli

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseArgs(t *testing.T) {
	Convey("with an empty slice", t, func() {
		result, err := ParseArgs([]string{})
		So(err, ShouldBeNil)
		So(result, ShouldResemble, Options{})
	})
	Convey("with a single filename", t, func() {
		result, err := ParseArgs([]string{"gim", "file.txt"})
		So(err, ShouldBeNil)
		So(result, ShouldResemble, Options{FilesToOpen: []string{"file.txt"}})
	})
	Convey("with multiple files", t, func() {
		result, err := ParseArgs([]string{"gim", "file1.txt", "file2.txt", "file3.txt"})
		So(err, ShouldBeNil)
		So(result, ShouldResemble, Options{FilesToOpen: []string{"file1.txt", "file2.txt", "file3.txt"}})
	})
	Convey("with --help", t, func() {
		result, err := ParseArgs([]string{"gim", "--help"})
		So(err, ShouldBeNil)
		So(result, ShouldResemble, Options{Help: true})
	})
	Convey("with invalid argument", t, func() {
		result, err := ParseArgs([]string{"gim", "--badarg"})
		So(err, ShouldNotBeNil)
		So(result, ShouldResemble, Options{})
	})
}
