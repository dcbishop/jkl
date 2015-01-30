package main

import (
	"testing"
	"time"

	"github.com/dcbishop/jkl/service"
	. "github.com/smartystreets/goconvey/convey"
)

func fakeApp() App {
	fd := NewFakeDriver()
	tui := NewTerminalUI(&fd)
	fs := GetCustomTestFs(fakeFileSystem)
	app := NewApp(fs)
	app.LoadOptions(SetUI(&tui))
	return app
}

func TestNew(t *testing.T) {
	Convey("app.New() returns basic App", t, func() {
		_ = fakeApp()
	})
}

func TestRunStop(t *testing.T) {
	Convey("given a new app", t, func() {
		app := fakeApp()
		So(app.Running(), ShouldBeFalse)

		Convey("app.Run() should start the app in a reasonable time", func() {
			go app.Run()
			So(service.WaitUntilRunning(&app, time.Second), ShouldBeNil)
			So(app.Running(), ShouldBeTrue)

			Convey("app.Stop() should terminate in a reasonable time", func() {
				go app.Stop()
				So(service.WaitUntilStopped(&app, time.Second), ShouldBeNil)
				So(app.Running(), ShouldBeFalse)
			})

			Convey("app.Run() on a running app shouldn't panic", func() {
				So(app.Run, ShouldNotPanic)
			})
		})

		Convey("app.Stop() without having called app.Run() shouldn't explode", func() {
			app.Stop()
			So(app.Running(), ShouldBeFalse)
		})
	})
}

func TestLoadOptions(t *testing.T) {
	app := fakeApp()
	Convey("accepts options", t, func() {
		options := Options{}

		app.LoadOptions(options...)
	})

	Convey("with 2 filenames given", t, func() {
		options, err := ParseArgs([]string{"jkl", "fakefile.txt", "fakefile2.txt"})
		So(err, ShouldBeNil)

		app.LoadOptions(options...)
		So(len(app.editor.Buffers()), ShouldEqual, 2)
	})
}
