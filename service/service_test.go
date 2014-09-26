package service_test

import (
	"testing"

	"github.com/dcbishop/gim/service"
	"github.com/dcbishop/gim/ui"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWaitUntilRunning(t *testing.T) {
	Convey("On a basic service", t, func() {
		ui := ui.TermboxUI{}
		So(ui.Running(), ShouldBeFalse)

		Convey("Timeout should return error", func() {
			err := service.WaitUntilRunning(&ui, 1)
			So(err, ShouldNotBeNil)
		})

		Convey("Should proceed when running state is met", func() {
			err := service.WaitUntilStopped(&ui, 1)
			So(err, ShouldBeNil)
		})
	})
}

func TestSetRunning(t *testing.T) {
	Convey("A blank state", t, func() {
		var state service.State
		Convey("Setting running should work", func() {
			So(state.Running(), ShouldBeFalse)
			err := state.SetRunning()

			So(state.Running(), ShouldBeTrue)
			So(err, ShouldBeNil)

			Convey("Unsetting running should work", func() {
				err := state.SetStopped()
				So(state.Running(), ShouldBeFalse)
				So(err, ShouldBeNil)
			})

			Convey("setting running when already running should return an error", func() {
				err := state.SetRunning()
				So(err, ShouldNotBeNil)
				So(state.Running(), ShouldBeTrue)
			})

		})
		Convey("Unsetting running when not running should return an error", func() {
			err := state.SetStopped()
			So(err, ShouldNotBeNil)
			So(state.Running(), ShouldBeFalse)
		})
	})
}
