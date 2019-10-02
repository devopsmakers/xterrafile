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
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"

	xt "github.com/devopsmakers/xterrafile/pkg/xterrafile"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the modules in your Terrafile",
	Run: func(cmd *cobra.Command, args []string) {
		jww.WARN.Printf("Removing all modules in %s", VendorDir)

		os.RemoveAll(VendorDir)
		os.MkdirAll(VendorDir, os.ModePerm)

		var wg sync.WaitGroup
		wg.Add(len(Config))

		for moduleName, moduleMeta := range Config {
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
		getWithGoGetter(moduleName, source, version, directory)
	default:
		getWithGoGetter(moduleName, moduleSource, moduleVersion, directory)
	}

	// If we have a path specified, let's extract it (move and copy stuff).
	if len(modulePath) > 0 {
		tmpDirectory := directory + ".tmp"
		pathWanted := path.Join(tmpDirectory, modulePath)

		err := os.Rename(directory, tmpDirectory)
		xt.CheckIfError(moduleName, err)

		err = copy.Copy(pathWanted, directory)
		xt.CheckIfError(moduleName, err)
		os.RemoveAll(tmpDirectory)
	}
	// Cleanup .git directory
	os.RemoveAll(path.Join(directory, ".git"))
}

// Handle modules from other sources to reflect:
// https://www.terraform.io/docs/modules/sources.html
//
// HEAVILY inpired by Terraform's internal getter / module_install code:
// https://github.com/hashicorp/terraform/blob/master/internal/initwd/getter.go
// https://github.com/hashicorp/terraform/blob/master/internal/initwd/module_install.go

var goGetterDetectors = []getter.Detector{
	new(getter.GitHubDetector),
	new(getter.BitBucketDetector),
	new(getter.GCSDetector),
	new(getter.S3Detector),
	new(getter.FileDetector),
}

var goGetterNoDetectors = []getter.Detector{}

var goGetterDecompressors = map[string]getter.Decompressor{
	"bz2": new(getter.Bzip2Decompressor),
	"gz":  new(getter.GzipDecompressor),
	"xz":  new(getter.XzDecompressor),
	"zip": new(getter.ZipDecompressor),

	"tar.bz2":  new(getter.TarBzip2Decompressor),
	"tar.tbz2": new(getter.TarBzip2Decompressor),

	"tar.gz": new(getter.TarGzipDecompressor),
	"tgz":    new(getter.TarGzipDecompressor),

	"tar.xz": new(getter.TarXzDecompressor),
	"txz":    new(getter.TarXzDecompressor),
}

var goGetterGetters = map[string]getter.Getter{
	"file":  new(getter.FileGetter),
	"gcs":   new(getter.GCSGetter),
	"git":   new(getter.GitGetter),
	"hg":    new(getter.HgGetter),
	"s3":    new(getter.S3Getter),
	"http":  getterHTTPGetter,
	"https": getterHTTPGetter,
}

var getterHTTPClient = cleanhttp.DefaultClient()

var getterHTTPGetter = &getter.HttpGetter{
	Client: getterHTTPClient,
	Netrc:  true,
}

func getWithGoGetter(name string, source string, version string, directory string) {

	// Fixup potential URLs for Github Detector
	if xt.IContains(source, ".git") {
		source = strings.Replace(source, "https://github.com/", "github.com/", 1)
	}

	moduleSource, err := getter.Detect(source, directory, getter.Detectors)
	xt.CheckIfError(name, err)

	jww.DEBUG.Printf("[%s] Detected real source: %s", name, moduleSource)

	realModuleSource, err := url.Parse(moduleSource)
	xt.CheckIfError(name, err)

	if len(version) > 0 {
		qParams := realModuleSource.Query()
		qParams.Set("ref", version)
		realModuleSource.RawQuery = qParams.Encode()
	}

	jww.INFO.Printf("[%s] Fetching %s", name, realModuleSource.String())
	client := getter.Client{
		Src: realModuleSource.String(),
		Dst: directory,
		Pwd: directory,

		Mode: getter.ClientModeDir,

		Detectors:     goGetterNoDetectors, // we already did detection above
		Decompressors: goGetterDecompressors,
		Getters:       goGetterGetters,
	}
	err = client.Get()
	xt.CheckIfError(name, err)
}
