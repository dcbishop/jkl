package ui

import (
	"testing"
	"time"

	"github.com/dcbishop/gim/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRunStop(t *testing.T) {
	Convey("ui.Run() should be terminated by ui.Stop() and ui.Running() should give the correct status", t, func() {
		ui := TermboxUI{}

		go ui.Run()

		if service.WaitUntilRunning(&ui, time.Second) != nil {
			panic("Failed to start UI.")
		}
		So(ui.Running(), ShouldBeTrue)

		ui.Stop()

		if service.WaitUntilStopped(&ui, time.Second) != nil {
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

		if service.WaitUntilRunning(&ui, time.Second) != nil {
			panic("Failed to start UI.")
		}

		So(ui.Running(), ShouldBeTrue)
		So(ui.Run, ShouldPanic)
		So(ui.Running(), ShouldBeTrue)
		ui.Stop()

		if service.WaitUntilStopped(&ui, time.Second) != nil {
			panic("Failed to start UI.")
		}
	})
}
