// Copyright © 2019 Tim Birkett <tim.birkett@devopsmakers.com>
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
	"fmt"
	"os"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type module struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

var cfgFile string
var VendorDir string
var Config map[string]module

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xterrafile",
	Short: "Manage any vendored modules with a YAML file.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	if (len(os.Args) > 1) && (os.Args[1] != "version") {
	  cobra.OnInitialize(initConfig)
	  rootCmd.PersistentFlags().StringVarP(&cfgFile, "file", "f", "Terrafile", "config file (default is: Terrafile)")
	  rootCmd.PersistentFlags().StringVarP(&VendorDir, "directory", "d", "vendor/xterrafile", "directory to download modules to")
  }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	yamlFile, err := ioutil.ReadFile(cfgFile)
  CheckIfError(err)

	err = yaml.Unmarshal(yamlFile, &Config)
	CheckIfError(err)
}
