package sensupluginschrony

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConvDate(t *testing.T) {
	var date string
	var time int64

	Convey("When checking the time conversion", t, func() {

		Convey("If the date is Thu Mar 17 18:09:31 2016", func() {
			date = "Thu Mar 17 18:09:31 2016"
			time = convDate(date)

			Convey("The time should be 1458238171", func() {
				So(time, ShouldEqual, 1458238171)
			})

			Convey("The time should not be 1458238131", func() {
				So(time, ShouldNotEqual, 1458238131)
			})
		})
	})
}

func TestOverIntThreshold(t *testing.T) {
	var num1 int64
	var num2 int64
	var alarm bool

	Convey("When the threshold is exceeded", t, func() {
		num1 = 20
		num2 = 2
		alarm = overIntThreshold(num1, num2)

		Convey("20 should be greater than 2 and the alarm should be true", func() {
			So(alarm, ShouldBeTrue)
		})
	})

	Convey("When the threshold is not exceeded", t, func() {
		num1 = 2
		num2 = 20
		alarm = overIntThreshold(num1, num2)

		Convey("2 should be less than 20 and the alarm should not be true", func() {
			So(alarm, ShouldBeFalse)
		})
	})
}

func TestOverFloatThreshold(t *testing.T) {
	var num1 float64
	var num2 float64
	var alarm bool

	Convey("When the threshold is exceeded", t, func() {
		num1 = 20.0
		num2 = 2.0
		alarm = overFloatThreshold(num1, num2)

		Convey("20.0 should be greater than 2.0 and the alarm should be true", func() {
			So(alarm, ShouldBeTrue)
		})
	})

	Convey("When the threshold is not exceeded", t, func() {
		num1 = 2.0
		num2 = 20.0
		alarm = overFloatThreshold(num1, num2)

		Convey("2.0 should be less than 20.0 and the alarm should not be true", func() {
			So(alarm, ShouldBeFalse)
		})
	})
}

func TestCleanVal(t *testing.T) {
	// var debug bool
	var c string
	var out string

	Convey("When debug is set to true", t, func() {
		c = "  ld50 caffeine"
		debug = true
		out = cleanVal(c)

		Convey("The output should be 'ld50 caffeine'", func() {
			So(out, ShouldEqual, "ld50 caffeine")
		})

		Convey("The output should not be '  ld50 caffeine'", func() {
			So(out, ShouldNotEqual, "  ld50 caffeine")
		})
	})

	Convey("When debug is set to false", t, func() {
		debug = false

		Convey("When the value input is ' 132.163.4.102 (time-b.timefreq.bldrdoc.gov)'", func() {
			c = " 132.163.4.102 (time-b.timefreq.bldrdoc.gov)"
			out = cleanVal(c)

			Convey("The output should be 'time-b.timefreq.bldrdoc.gov'", func() {
				So(out, ShouldEqual, "time-b.timefreq.bldrdoc.gov")
			})

			Convey("The output should not be anything else", func() {
				So(out, ShouldNotEqual, "(time-b.timefreq.bldrdoc.gov)")
				So(out, ShouldNotEqual, "132.163.4.102 (time-b.timefreq.bldrdoc.gov)")
				So(out, ShouldNotEqual, "132.163.4.102")
			})
		})

		Convey("When the value input is ' 0.003938113 seconds fast of NTP time'", func() {
			c = " 0.003938113 seconds fast of NTP time"
			out = cleanVal(c)

			Convey("The output should be '0.003938113'", func() {
				So(out, ShouldEqual, "0.003938113")
			})

			Convey("The output should not be anything else", func() {
				So(out, ShouldNotEqual, " 0.003938113 seconds fast of NTP time")
				So(out, ShouldNotEqual, "0.003938113 seconds fast of NTP time")
			})
		})

		Convey("When the value input is ' -0.063844 seconds'", func() {
			c = " -0.063844 seconds"
			out = cleanVal(c)

			Convey("The output should be '-0.063844'", func() {
				So(out, ShouldEqual, "-0.063844")
			})

			Convey("The output should not be anything else", func() {
				So(out, ShouldNotEqual, "0.063844 seconds")
				So(out, ShouldNotEqual, "-0.063844 seconds")
				So(out, ShouldNotEqual, "0.063844")
			})
		})
	})
}

func TestCreateMap(t *testing.T) {

	Convey("When creating a map", t, func() {
		var m = make(map[string]string)

		Convey("If the string is 'System time     : 0.003938113 seconds fast of NTP time'", func() {
			input := "System time     : 0.066877037 seconds slow of NTP time"
			m = createMap(input)

			Convey("The value of 'System time' should be '0.066877037'", func() {
				So(m["System time"], ShouldEqual, "0.066877037")
			})

			Convey("The key should not be anything else", func() {
				So(m["System time     "], ShouldNotEqual, "0.000291823 seconds slow of NTP time")
				So(m["System time     "], ShouldNotEqual, "0.000291823")
				So(m["System time"], ShouldNotEqual, "0.000291823 seconds slow of NTP time")
			})
		})
	})
}

func TestTimeDeviation(t *testing.T) {
	Convey("When testing the difference of two times", t, func() {
		var curTime int64
		var checkTime int64
		var threshold int64
		var out bool

		Convey("When '1458528828' minus '1458521701' is over the threshold of 300 seconds", func() {
			curTime = 1458528828
			checkTime = 1458521701
			threshold = 300

			out = timeDeviation(checkTime, curTime, threshold)
			Convey("The output should be true", func() {
				So(out, ShouldBeTrue)
			})
			Convey("The output should not be false", func() {
				So(out, ShouldNotEqual, false)
			})
		})

		Convey("When '1458528828' minus '1458528820' is under the threshold of 300 seconds", func() {
			curTime = 1458528828
			checkTime = 1458528820
			threshold = 300

			out = timeDeviation(checkTime, curTime, threshold)
			Convey("The output should be false", func() {
				So(out, ShouldBeFalse)
			})
			Convey("The output should not be true", func() {
				So(out, ShouldNotEqual, true)
			})
		})
	})
}
