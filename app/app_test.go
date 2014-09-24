package app

import (
	"testing"

	"github.com/dcbishop/gim/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("app.New() returns basic App", t, func() {
		app := New()
		So(app, ShouldResemble, App{})
	})
}

func TestLoadOptions(t *testing.T) {
	app := New()
	Convey("with default Options returns basic App", t, func() {
		options := cli.Options{}

		app.LoadOptions(options)

		So(app, ShouldResemble, App{})
	})
}
