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
	"os"
	"os/exec"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/yieldbot/sensuplugin/sensuutil"
)

var warnThreshold int64
var critThreshold int64
var checkKey string

var condition string
var msg string

var syslogLog = logging.MustGetLogger("chrony")
var stderrLog = logging.MustGetLogger("chrony")

var checkChronyStatsCmd = &cobra.Command{
	Use:   "checkChronyStats",
	Short: "Check various values in chrony to ensure all is well",
	Long: `This will use 'chronyc tracking' to build a map of keys allowing the
  user to check against any of the values to ensure they are within tolerated
  limits for their environment.

  Currently the following values can be checked:
  - Refernce ID
  - Stratum
  - Reference Time
  - Last Offset
  - RMS Offset`,

	Run: func(sensupluginschrony *cobra.Command, args []string) {

		syslogBackend, _ := logging.NewSyslogBackend("checkChronyStats")
		stderrBackend := logging.NewLogBackend(os.Stderr, "checkChronyStats", 0)
		syslogBackendFormatter := logging.NewBackendFormatter(syslogBackend, sensuutil.SyslogFormat)
		stderrBackendFormatter := logging.NewBackendFormatter(stderrBackend, sensuutil.StderrFormat)
		logging.SetBackend(syslogBackendFormatter)
		logging.SetBackend(stderrBackendFormatter)

		chronyStats := exec.Command("chronyc", "tracking")

		out, err := chronyStats.Output()
		if err != nil {
			syslogLog.Error("err")
			os.Exit(129)
		}

		chronyStats.Start()
		data := createMap(string(out))

		if debug {
			for k, v := range data {
				stderrLog.Debug("Key: ", k, "Current value: ", v)
			}
			sensuutil.Exit("Debug")
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
			sensuutil.Exit("ok")
		case "warning":
			sensuutil.Exit("warning", msg)
		case "critical":
			sensuutil.Exit("critical", msg)
		}
	},
}

func init() {
	RootCmd.AddCommand(checkChronyStatsCmd)

	checkChronyStatsCmd.Flags().Int64VarP(&warnThreshold, "warn", "", 4, "the alert warning threshold")
	checkChronyStatsCmd.Flags().Int64VarP(&critThreshold, "crit", "", 8, "the alert critical threshold")
	checkChronyStatsCmd.Flags().StringVarP(&checkKey, "checkKey", "", "", "the key to check")
}

/*

ReferenceID is a straight shot, it is either critical or not
Stratum is the number of hops away, critical and warning are that number
ReferenceTime is the time the last measurement from a source was processed, critical and warning values represent the number of seconds diviation

*/
