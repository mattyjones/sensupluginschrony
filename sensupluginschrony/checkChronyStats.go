// Copyright © 2016 Yieldbot <devops@yieldbot.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sensupluginschrony

import (
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yieldbot/sensuplugin/sensuutil"
)

var warnThreshold int64
var critThreshold int64

var condition string
var msg string

var checkChronyStatsCmd = &cobra.Command{
	Use:   "checkChronyStats",
	Short: "Check various values in chrony to ensure all is well",
	Long: `This will use 'chronyc tracking' to build a map of keys allowing the user to check against any of the values to ensure they are within tolerated limits for their environment.

  Currently the following values can be checked:
  - Refernce ID
  - Stratum
  - Reference Time
  - Last Offset
  - RMS Offset`,

	Run: func(sensupluginschrony *cobra.Command, args []string) {

		chronyStats := exec.Command("chronyc", "tracking")

		out, err := chronyStats.Output()
		if err != nil {
			syslogLog.WithFields(logrus.Fields{
				"check":   "sensupluginscrony",
				"client":  host,
				"version": "foo",
				"error":   err,
				"output":  out,
			}).Error(`ChronyStats output is not valid`)
			sensuutil.Exit("RUNTIMEERROR")
		}

		chronyStats.Start()
		data := createMap(string(out))

		if debug {
			for k, v := range data {
				syslogLog.WithFields(logrus.Fields{
					"check":         "sensupluginscrony",
					"client":        host,
					"version":       "foo",
					"key":           k,
					"Current value": v,
				}).Info()
			}
			sensuutil.Exit("DEBUG")
		}

		switch checkKey {
		case "ReferenceID":
			condition, msg = checkLocalChrony(data["Reference ID"])
		case "Stratum":
			condition, msg = checkStratum(data["Stratum"], warnThreshold, critThreshold)
		case "ReferenceTime":
			condition, msg = checkRefTime(convDate(data["Ref time (UTC)"]), warnThreshold, critThreshold)
		case "LastOffset":
			condition, msg = checkOffset(data["Last offset"], warnThreshold, critThreshold)
		case "RMSOffset":
			condition, msg = checkOffset(data["RMS offset"], warnThreshold, critThreshold)
		}

		switch condition {
		case "ok":
			sensuutil.Exit(condition)
		case "warning":
			sensuutil.Exit(condition, msg)
		case "critical":
			sensuutil.Exit(condition, msg)
		}
	},
}

func init() {
	RootCmd.AddCommand(checkChronyStatsCmd)

	checkChronyStatsCmd.Flags().Int64VarP(&warnThreshold, "warn", "", 0, "the alert warning threshold")
	checkChronyStatsCmd.Flags().Int64VarP(&critThreshold, "crit", "", 0, "the alert critical threshold")
}

/*

ReferenceID is a straight shot, it is either critical or not
Stratum is the number of hops away, critical and warning are that number
ReferenceTime is the time the last measurement from a source was processed, critical and warning values represent the number of seconds diviation

*/
