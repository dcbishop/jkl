package main

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
	Convey("with invalid argument", t, func() {
		result, err := ParseArgs([]string{"jkl", "--badarg"})
		So(err, ShouldNotBeNil)
		So(result, ShouldResemble, Options{})
	})
}
