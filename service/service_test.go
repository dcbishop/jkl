package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type FakeService struct {
	state State
}

func (service *FakeService) Run() {
	if service.state.SetRunning() != nil {
		panic("Service already running")
	}
}

func (service *FakeService) Stop() {
	service.state.SetStopped()
}

func (service *FakeService) Running() bool {
	return service.state.Running()
}

func TestWaitUntilRunning(t *testing.T) {
	Convey("On a basic service", t, func() {
		service := FakeService{}
		So(service.Running(), ShouldBeFalse)

		Convey("Timeout should return error", func() {
			err := WaitUntilRunning(&service, 1)
			So(err, ShouldNotBeNil)
		})

		Convey("Should proceed when running state is met", func() {
			err := WaitUntilStopped(&service, 1)
			So(err, ShouldBeNil)
		})
	})
}

func TestSetRunning(t *testing.T) {
	Convey("A blank state", t, func() {
		var state State
		So(state.Running(), ShouldBeFalse)

		Convey("Setting running should work", func() {
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
