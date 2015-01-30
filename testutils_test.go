package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringToRunes(t *testing.T) {
	Convey("Basic 3x3 unicode box", t, func() {
		expected := [][]rune{
			{'╔', '═', '╗'},
			{'║', 0, '║'},
			{'╚', '═', '╝'},
		}

		runegrid := StringToRunes(UnicodeBox, ' ')
		So(len(runegrid), ShouldEqual, 3)
		So(runegrid, ShouldResemble, expected)
	})
}
