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

package xterrafile

import (
	"net/url"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"
	jww "github.com/spf13/jwalterweatherman"
)

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

// GetWithGoGetter downloads objects using go-getter
func GetWithGoGetter(name string, source string, version string, directory string) {

	// Fixup potential URLs for Github Detector
	if IContains(source, ".git") {
		source = strings.Replace(source, "https://github.com/", "github.com/", 1)
	}

	moduleSource, err := getter.Detect(source, directory, getter.Detectors)
	CheckIfError(name, err)

	jww.DEBUG.Printf("[%s] Detected real source: %s", name, moduleSource)

	realModuleSource, err := url.Parse(moduleSource)
	CheckIfError(name, err)

	qParams := realModuleSource.Query()

	if len(version) > 0 && len(qParams) == 0 {
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
	CheckIfError(name, err)
}
