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
	"regexp"
	"strings"
	"time"
)

func createMap(out string) map[string]string {
	m := make(map[string]string)
	var key string
	var c string

	j := strings.Split(out, "\n")
	for _, k := range j {
		l := strings.SplitN(k, ":", 2)
		for i, n := range l {
			c = string(n)
			if i == 0 {
				if !debug {
					key = strings.TrimSpace(c)
				} else {
					key = c
				}
			} else {
				m[key] = cleanVal(c)
			}
		}
	}
	return m
}

func cleanVal(c string) string {
	reNum := regexp.MustCompile(`^\s+-?[0-9]+`)
	reIP := regexp.MustCompile(`^\s+[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
	out := c
	if !debug {
		if reIP.MatchString(c) {
			h := strings.Split(c, " ")
			for i, s := range h {
				if i == 2 {
					r := strings.NewReplacer("(", "", ")", "")
					out = string(r.Replace(s))
				}
			}
		} else if reNum.MatchString(c) {
			h := strings.Split(c, " ")
			for i, s := range h {
				if i == 1 {
					out = string(s)
				}
			}
		}
	}
	return strings.TrimSpace(out)
}

func timeDeviation(val int64, curVal int64, threshold int64) bool {
	if (curVal - val) > threshold {
		return true
	}
	return false
}

func convDate(d string) int64 {
	e, _ := time.Parse("Mon Jan _2 15:04:05 2006", d)
	return e.Unix()
}

func overIntThreshold(num1 int64, num2 int64) bool {
	if num1 >= num2 {
		return true
	}
	return false
}

func overFloatThreshold(num1 float64, num2 float64) bool {
	if num1 >= num2 {
		return true
	}
	return false
}
