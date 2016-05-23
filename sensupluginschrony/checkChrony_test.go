package sensupluginschrony

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckLocalChrony(t *testing.T) {
	var condition string
	// var msg string

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

func TestCheckStratum(t *testing.T) {
	var alarm string
	var curVal string
	var warnThreshold int64
	var critThreshold int64

	Convey("When checking the Stratum", t, func() {
		Convey("If the Stratum is at 12 and the critical level is 10 and the warning level is 5", func() {
			curVal = "12"
			warnThreshold = 5
			critThreshold = 10
			alarm, msg = checkStratum(curVal, warnThreshold, critThreshold)

			Convey("The alarm should be true and the message should be 'You are more than the max number of hops'", func() {
				So(alarm, ShouldEqual, "critical")
				So(msg, ShouldEqual, "You are more than the max number of hops")
			})
			Convey("The alarm should not be warning and the message should not be 'You are nearing the max number of hops'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are nearing the max number of hops")
			})
			Convey("The alarm should not be 'ok' and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
		})

		Convey("If the Stratum is at 6 and the critical level is 10 and the warning level is 5", func() {
			curVal = "6"
			warnThreshold = 5
			critThreshold = 10
			alarm, msg = checkStratum(curVal, warnThreshold, critThreshold)

			Convey("The alarm should be 'critical' and the message should not be 'You are more than the max number of hops'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are more than the max number of hops")
			})
			Convey("The alarm should be warning and the message should be 'You are nearing the max number of hops'", func() {
				So(alarm, ShouldEqual, "warning")
				So(msg, ShouldEqual, "You are nearing the max number of hops")
			})
			Convey("The alarm should not be 'ok' and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
		})

		Convey("If the Stratum is at 4 and the critical level is 8 and the warning level is 5", func() {
			curVal = "4"
			warnThreshold = 5
			critThreshold = 10
			alarm, msg = checkStratum(curVal, warnThreshold, critThreshold)

			Convey("The alarm should not be 'critical' and the message should not be 'You are more than the max number of hops'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are more than the max number of hops")
			})
			Convey("The alarm should not be 'warning' and the message should not be 'You are nearing the max number of hops'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are nearing the max number of hops")
			})
			Convey("The alarm should be 'ok' and the message should be empty", func() {
				So(alarm, ShouldEqual, "ok")
				So(msg, ShouldBeEmpty)
			})
		})
	})
}

func TestCheckRefTime(t *testing.T) {
	var warnThreshold int64
	var critThreshold int64
	var alarm string
	var msg string
	var curVal int64

	Convey("When checking the Reference Time", t, func() {
		Convey("If the reference time is at 1000 and we compare against the current time and the critical threshold is 300s", func() {
			curVal = 1000
			warnThreshold = 60
			critThreshold = 300
			alarm, msg = checkRefTime(curVal, warnThreshold, critThreshold)

			Convey("The alarm should be critical and the message should be 'You are over the max allowed deviation'", func() {
				So(alarm, ShouldEqual, "critical")
				So(msg, ShouldEqual, "You are over the max allowed deviation")
			})
			Convey("The alarm not should be warning and the message should not be 'You are nearing the max allowed deviation'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are nearing the max allowed deviation")
			})
			Convey("The alarm should not be 'ok' and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
		})

		Convey("If the reference time is the current time minus 100s and we compare against the current time and the warning threshold is 60s", func() {

			curVal = time.Now().UTC().Unix() - 100
			warnThreshold = 60
			critThreshold = 300
			alarm, msg = checkRefTime(curVal, warnThreshold, critThreshold)

			Convey("The alarm should not be critical and the message should not be 'You are over the max allowed deviation'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are over the max allowed deviation")
			})
			Convey("The alarm should be warning and the message should be 'You are nearing the max allowed deviation'", func() {
				So(alarm, ShouldEqual, "warning")
				So(msg, ShouldEqual, "You are nearing the max allowed deviation")
			})
			Convey("The alarm should not be 'ok' and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
		})

		Convey("If the reference time is the current time and we compare against the current time and the critical threshold is 300s", func() {

			curVal = time.Now().UTC().Unix()
			warnThreshold = 100
			critThreshold = 300
			alarm, msg = checkRefTime(curVal, warnThreshold, critThreshold)

			Convey("The alarm should not be critical and the message should not be 'You are over the max allowed deviation'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are over the max allowed deviation")
			})
			Convey("The alarm not should be warning and the message should not be 'You are nearing the max allowed deviation'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are nearing the max allowed deviation")
			})
			Convey("The alarm should be 'ok' and the message should be empty", func() {
				So(alarm, ShouldEqual, "ok")
				So(msg, ShouldBeEmpty)
			})
		})
	})
}

func TestCheckOffset(t *testing.T) {
	var warnThreshold int64
	var critThreshold int64
	var alarm string
	var msg string
	var offset string

	Convey("When checking the Offset", t, func() {
		Convey("If the offset is 3.14, the warning threshold is 4, and the critical threshold is 6", func() {
			offset = "3.14"
			warnThreshold = 4
			critThreshold = 6
			alarm, msg = checkOffset(offset, warnThreshold, critThreshold)

			Convey("The alarm should be ok and the message should be empty", func() {
				So(alarm, ShouldEqual, "ok")
				So(msg, ShouldBeEmpty)
			})
			Convey("The alarm not should be warning and the message should not be 'You are over the warning threshold'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are over the warning threshold")
			})
			Convey("The alarm should not be 'critical' and the message should not be 'You are over the critical threshold'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are over the critical threshold")
			})
		})

		Convey("If the offset is 4.14, the warning threshold is 4, and the critical threshold is 6", func() {

			offset = "4.14"
			warnThreshold = 4
			critThreshold = 6
			alarm, msg = checkOffset(offset, warnThreshold, critThreshold)

			Convey("The alarm should be not be ok and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
			Convey("The alarm should be warning and the message should be 'You are over the warning threshold'", func() {
				So(alarm, ShouldEqual, "warning")
				So(msg, ShouldEqual, "You are over the warning threshold")
			})
			Convey("The alarm should not be 'critical' and the message should not be 'You are over the critical threshold'", func() {
				So(alarm, ShouldNotEqual, "critical")
				So(msg, ShouldNotEqual, "You are over the critical threshold")
			})
		})

		Convey("If the offset is 6.14, the warning threshold is 4, and the critical threshold is 6", func() {

			offset = "6.14"
			warnThreshold = 4
			critThreshold = 6
			alarm, msg = checkOffset(offset, warnThreshold, critThreshold)

			Convey("The alarm should be not be ok and the message should not be empty", func() {
				So(alarm, ShouldNotEqual, "ok")
				So(msg, ShouldNotBeEmpty)
			})
			Convey("The alarm should not be warning and the message should not be 'You are over the warning threshold'", func() {
				So(alarm, ShouldNotEqual, "warning")
				So(msg, ShouldNotEqual, "You are over the warning threshold")
			})
			Convey("The alarm should be 'critical' and the message should be 'You are over the critical threshold'", func() {
				So(alarm, ShouldEqual, "critical")
				So(msg, ShouldEqual, "You are over the critical threshold")
			})
		})
	})
}
