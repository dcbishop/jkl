package ui

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWaitUntilServiceRunning(t *testing.T) {
	Convey("Timeout should return error", t, func() {
		ui := TermboxUI{}
		err := WaitUntilServiceRunning(&ui, true, 1)
		So(err, ShouldNotBeNil)
	})
	Convey("Should proceed when running state is met", t, func() {
		ui := TermboxUI{}
		err := WaitUntilServiceRunning(&ui, false, 1)
		So(err, ShouldBeNil)
	})
}

func TestRunStop(t *testing.T) {
	Convey("ui.Run() should be terminated by ui.Stop() and ui.Running() should give the correct status", t, func() {
		ui := TermboxUI{}

		go ui.Run()

		if WaitUntilServiceRunning(&ui, true, time.Second) != nil {
			panic("Failed to start UI.")
		}
		So(ui.Running(), ShouldBeTrue)

		ui.Stop()

		if WaitUntilServiceRunning(&ui, false, time.Second) != nil {
			panic("Failed to stop UI.")
		}
		So(ui.Running(), ShouldBeFalse)

	})
	Convey("ui.Stop() at the same time as ui.Run()", t, func() {
		ui := TermboxUI{}
		go ui.Stop()
		ui.Run()
		So(ui.Running(), ShouldBeFalse)
	})
	Convey("ui.Stop() without having called ui.Run() shouldn't explode", t, func() {
		ui := TermboxUI{}
		ui.Stop()
		ui.Stop()
		ui.Stop()
		So(ui.Running(), ShouldBeFalse)
	})
	Convey("ui.Run() when already running should panic", t, func() {
		ui := TermboxUI{}
		go ui.Run()

		if WaitUntilServiceRunning(&ui, true, time.Second) != nil {
			panic("Failed to start UI.")
		}

		So(ui.Running(), ShouldBeTrue)
		So(ui.Run, ShouldPanic)
		So(ui.Running(), ShouldBeTrue)
		ui.Stop()

		if WaitUntilServiceRunning(&ui, false, time.Second) != nil {
			panic("Failed to start UI.")
		}
	})
}
