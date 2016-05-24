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
	"fmt"
	"log/syslog"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yieldbot/sensuplugin/sensuutil"
)

var cfgFile string // used for configuration via Viper

var host string // get the hostname for logging

// Create a set of logging instances. I have two because they are configured
// differently via levels and format. Stderr should be configured for ascii
// as that is more human readable but the syslog logger should be in json format
// to make it more easily consumable via automated processes or third-party
// tools.
// var stderrLog = logrus.Logger{
// 	Out: os.Stderr,
// }
// var stdoutLog = logrus.Logger{
// 	Out:   os.Stdout,
// 	Level: logrus.DebugLevel,
// }
var syslogLog = logrus.New()

var checkKey string

var debug bool

// RootCmd is the toplevel entry point
var RootCmd = &cobra.Command{
	Use:   "sensupluginschrony",
	Short: "kk",
	Long: `This plugin uses 'chronyc tracking' to build a map of values and then compares the given value against a warning and critical threshold supplied by the user either via command line or via configuration file

  Currently the following values can be checked:
  - Refernce ID
  - Stratum
  - Reference Time
  - Last Offset
  - RMS Offset`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Setup logging for the package. Doing it here is much eaiser than in each
	// binary. If you want to overwrite it in a specific binary then feel free.
	// stderrLog.Out = os.Stderr
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		panic(err)
	}
	syslogLog.Hooks.Add(hook)
	// syslogLog.Formatter = new(logrus.JSONFormatter)

	// Set the hostname for use in logging within the package. Doing it here is
	// cleaner than in each binary but if you want to use some other method just
	// override the variable in the specific binary.
	host, err = os.Hostname()
	if err != nil {
		syslogLog.WithFields(logrus.Fields{
			"check":   "sensupluginschrony",
			"client":  "unknown",
			"version": "foo",
			"error":   err,
		}).Error(`Could not determine the hostname of this machine as reported by the kernel.`)
		sensuutil.Exit("GENERALGOLANGERROR")
	}

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/sensuplugins/conf.d/.sensupluginschrony.yaml)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debugging info")
	RootCmd.PersistentFlags().StringVarP(&checkKey, "checkKey", "", "", "the key to check")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("sensupluginschrony")
		viper.AddConfigPath("/etc/sensuplugins/conf.d")
	}

	if checkKey == "" {
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err == nil {
		} else {
			syslogLog.WithFields(logrus.Fields{
				"check":   "sensupluginscrony",
				"client":  host,
				"version": "foo",
				"error":   err,
				"cfgFile": cfgFile,
			}).Error(`Could not read in the configuration file.`)
		}
	}
}
