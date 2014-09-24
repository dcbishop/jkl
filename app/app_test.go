package app

import (
	"testing"

	"github.com/dcbishop/gim/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("New", t, func() {
		Convey("Default Options", func() {
			app := New(cli.Options{})
			So(app, ShouldResemble, App{})
		})
	})
}
