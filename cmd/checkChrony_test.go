package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckLocalChrony(t *testing.T) {
	var condition string
	var msg string

	Convey("When checking the Reference ID", t, func() {
		var RefID string

		Convey("If the RefID is a local IP", func() {
			RefID = "127.127.1.1"
			condition, msg = checkLocalChrony(RefID)

			Convey("The alert status should be critical and the message should not be empty", func() {
				So(condition, ShouldEqual, "critical")
				So(msg, ShouldEqual, "Chrony is synced locally")
			})

			Convey("The alert status should not be ok and the message should be empty", func() {
				So(condition, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
		})
		Convey("If the RefID is a remote IP", func() {
			RefID = "8.8.8.8"
			condition, msg = checkLocalChrony(RefID)

			Convey("The alert status should be ok and the message should be empty", func() {
				So(condition, ShouldEqual, "ok")
				So(msg, ShouldBeEmpty)
			})

			Convey("The alert status should not be critical and the message should not be empty", func() {
				So(condition, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "Chrony is synced locally")
			})
		})
	})
}
