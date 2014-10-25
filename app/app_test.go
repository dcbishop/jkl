package app

import (
	"testing"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/cli"
	"github.com/dcbishop/jkl/service"
	"github.com/dcbishop/jkl/ui"
	. "github.com/smartystreets/goconvey/convey"
)

var fakeFileName = "fakefile.txt"
var fakeFileContents = []byte(`Hello, this is a test`)
var fakeFileContents2 = []byte(`This is the 2nd file!`)

var fakeFileSystem = map[string][]byte{
	"fakefile.txt":  fakeFileContents,
	"fakefile2.txt": fakeFileContents,
}

var fakeFileAccessor = fileaccessor.Virtual{fakeFileSystem}

func fakeApp() App {
	fd := ui.NewFakeDriver()
	tui := ui.NewTerminalUI(&fd)
	app := New(fakeFileAccessor, &tui)
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

			Convey("app.Run() on a running app should panic", func() {
				So(app.Run, ShouldPanic)
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
		options := cli.Options{}

		app.LoadOptions(options)
	})

	Convey("with 2 filenames given", t, func() {
		options, err := cli.ParseArgs([]string{"jkl", "fakefile.txt", "fakefile2.txt"})
		So(err, ShouldBeNil)

		app.LoadOptions(options)
		So(len(app.editor.Buffers()), ShouldEqual, 2)
	})
}
