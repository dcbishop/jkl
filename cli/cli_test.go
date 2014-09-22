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
}
