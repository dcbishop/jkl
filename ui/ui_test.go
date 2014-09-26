package ui

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRunStop(t *testing.T) {
	Convey("ui.Run() should be terminated by an ui.Stop()", t, func() {
		ui := TermboxUI{}
		time.AfterFunc(time.Second/5, func() {
			panic("Failed to stop UI.")
		})
		go ui.Stop()
		ui.Run()
	})
	Convey("ui.Run() should be terminated by an ui.Stop() when blocking", t, func() {
		ui := TermboxUI{}
		time.AfterFunc(time.Second/10, func() {
			ui.Stop()
		})
		time.AfterFunc(time.Second/5, func() {
			panic("Failed to stop UI.")
		})
		ui.Run()
	})
	Convey("ui.Stop() without having called ui.Run() shouldn't explode", t, func() {
		ui := TermboxUI{}
		ui.Stop()
		ui.Stop()
		ui.Stop()
	})
}
