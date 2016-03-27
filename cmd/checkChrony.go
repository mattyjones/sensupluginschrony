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

package cmd

import (
	"strconv"
	"time"
)

func checkLocalChrony(RefID string) (string, string) {
	if RefID == "127.127.1.1" {
		msg := "Chrony is synced locally"
		return "critical", msg
	}
	return "ok", ""
}

func checkStratum(curVal string, warnThreshold int64, critThreshold int64) (string, string) {
	if val, err := strconv.ParseInt(curVal, 10, 32); err == nil {
		switch {
		case overIntThreshold(val, critThreshold):
			msg := "You are more than the max number of hops"
			return "critical", msg
		case overIntThreshold(val, warnThreshold):
			msg := "You are nearing the max number of hops"
			return "warning", msg
		}
	}
	return "ok", ""
}

func checkRefTime(curVal int64, warnThreshold int64, critThreshold int64) (string, string) {
	val := curVal
	t := time.Now().UTC().Unix()
	switch {
	case timeDeviation(val, t, critThreshold):
		msg := "You are over the max allowed deviation"
		return "critical", msg
	case timeDeviation(val, t, warnThreshold):
		msg := "You are nearing the max allowed deviation"
		return "warning", msg
	}
	return "ok", ""
}

func checkOffset(offset string, warnThreshold int64, critThreshold int64) (string, string) {
	curVal, _ := strconv.ParseFloat(offset, 64)

	switch {
	case overFloatThreshold(curVal, float64(critThreshold)):
		msg := "You are over the critical threshold"
		return "critical", msg
	case overFloatThreshold(curVal, float64(warnThreshold)):
		msg := "You are over the warning threshold"
		return "warning", msg
	}
	return "ok", ""
}
