package ui

import (
	"testing"
	"time"

	"github.com/dcbishop/fileaccessor"
	"github.com/dcbishop/jkl/editor"
	"github.com/dcbishop/jkl/runegrid"
	"github.com/dcbishop/jkl/service"
	"github.com/dcbishop/jkl/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFakeUI(t *testing.T) {
	Convey("A new FakeUI should be in a stopped state", t, func() {
		ui := FakeUI{}
		testUI(&ui)
	})
}

func TestTermbox(t *testing.T) {
	Convey("A new TerminalUI should be in a stopped state", t, func() {
		ui := NewTerminalUI()
		ui.Console = &FakeDriver{}
		testUI(&ui)
	})
}

func testUI(ui UI) {
	So(ui.Running(), ShouldBeFalse)

	Convey("Run() should set state to Running() in a reasonable time", func() {
		go ui.Run()
		if service.WaitUntilRunning(ui, time.Second) != nil {
			panic("Failed to start UI.")
		}
		So(ui.Running(), ShouldBeTrue)

		Convey("should be terminated by Stop() and Running() should give the correct status", func() {
			ui.Stop()

			if service.WaitUntilStopped(ui, time.Second) != nil {
				panic("Failed to stop UI.")
			}
			So(ui.Running(), ShouldBeFalse)

		})

		Convey("Run() when already running should panic", func() {
			So(ui.Running(), ShouldBeTrue)
			So(ui.Run, ShouldPanic)
			So(ui.Running(), ShouldBeTrue)
			ui.Stop()

			if service.WaitUntilStopped(ui, time.Second) != nil {
				panic("Failed to start UI.")
			}
		})

		Convey("Event() should now have a valid channel", func() {
			events := ui.Events()
			So(events, ShouldNotBeNil)
		})
	})

	Convey("Stop() at the same time as ui.Run()", func() {
		go ui.Stop()
		ui.Run()
		So(ui.Running(), ShouldBeFalse)
	})

	Convey("Stop() without having called Run() shouldn't explode", func() {
		ui.Stop()
		ui.Stop()
		ui.Stop()
		So(ui.Running(), ShouldBeFalse)
	})
}

func TestRedraw(t *testing.T) {
	Convey("Basic Editor.", t, func() {
		tui := NewTerminalUI()
		console := FakeDriver{}
		console.Width = 3
		console.Height = 3
		tui.Console = &console

		go tui.Run()
		service.WaitUntilRunning(&tui, time.Second)

		fa := fileaccessor.Virtual{map[string][]byte{
			"file.txt": {'!'},
		}}
		editor := editor.New(fa)
		editor.OpenFile("file.txt")

		expected := runegrid.StringToRuneGrid(testhelpers.Empty3x3, '.')
		So(console.Grid, ShouldResemble, expected)

		Convey("Render 3x3 with default settings", func() {
			UnicodeBox := `
			╔═╗
			║!║
			╚═╝
			`
			expected := runegrid.StringToRuneGrid(UnicodeBox, 0)
			tui.Redraw(&editor)

			So(console.Grid, ShouldResemble, expected)
			So(true, ShouldBeTrue)
		})
	})
}
