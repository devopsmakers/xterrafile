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
	"os"
	"path"
	"sync"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	getter "github.com/hashicorp/go-getter"

	xt "github.com/devopsmakers/xterrafile/pkg"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the modules in your Terrafile",
	Run: func(cmd *cobra.Command, args []string) {
		jww.WARN.Printf("Removing all modules in %s", VendorDir)

		_ = os.RemoveAll(VendorDir)
		_ = os.MkdirAll(VendorDir, os.ModePerm)

		var wg sync.WaitGroup
		wg.Add(len(Config.Modules))

		for moduleName, moduleMeta := range Config.Modules {
			go getModule(moduleName, moduleMeta, &wg)
		}
		wg.Wait()
	},
}

func init() {
	jww.SetStdoutThreshold(jww.LevelInfo)
	rootCmd.AddCommand(installCmd)
}

func getModule(moduleName string, moduleMeta module, wg *sync.WaitGroup) {
	defer wg.Done()

	moduleSource, modulePath := getter.SourceDirSubdir(moduleMeta.Source)

	moduleVersion := ""
	if len(moduleMeta.Version) > 0 {
		moduleVersion = moduleMeta.Version
	}

	if len(moduleMeta.Path) > 0 {
		modulePath = moduleMeta.Path
	}

	directory := path.Join(VendorDir, moduleName)

	switch {
	case xt.IsLocalSourceAddr(moduleSource):
		xt.CopyFile(moduleName, moduleSource, directory)
	case xt.IsRegistrySourceAddr(moduleSource):
		source, version := xt.GetRegistrySource(moduleName, moduleSource, moduleVersion, nil)
		xt.GetWithGoGetter(moduleName, source, version, directory)
	case xt.IsGitSourceAddr(moduleSource):
		source, version := xt.GetGitSource(moduleName, moduleSource, moduleVersion)
		xt.GetWithGoGetter(moduleName, source, version, directory)
	default:
		xt.GetWithGoGetter(moduleName, moduleSource, moduleVersion, directory)
	}

	// If we have a path specified, let's extract it (move and copy stuff).
	if len(modulePath) > 0 {
		tmpDirectory := directory + ".tmp"
		pathWanted := path.Join(tmpDirectory, modulePath)

		err := os.Rename(directory, tmpDirectory)
		xt.CheckIfError(moduleName, err)

		err = copy.Copy(pathWanted, directory)
		xt.CheckIfError(moduleName, err)
		_ = os.RemoveAll(tmpDirectory)
	}
	// Cleanup .git directory
	_ = os.RemoveAll(path.Join(directory, ".git"))
}
