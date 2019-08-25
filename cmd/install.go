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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"
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

	moduleSource, modulePath := splitAddrSubdir(moduleMeta.Source)

	moduleVersion := ""
	if len(moduleMeta.Version) > 0 {
		moduleVersion = moduleMeta.Version
	}

	if len(moduleMeta.Path) > 0 {
		modulePath = moduleMeta.Path
	}

	directory := path.Join(VendorDir, moduleName)

	switch {
	case isLocalSourceAddr(moduleSource):
		copyFile(moduleName, moduleSource, directory)
	case isRegistrySourceAddr(moduleSource):
		source, version := getRegistrySource(moduleName, moduleSource, moduleVersion)
		getWithGoGetter(moduleName, source, version, directory)
	default:
		getWithGoGetter(moduleName, moduleSource, moduleVersion, directory)
	}

	// If we have a path specified, let's extract it (move and copy stuff).
	if len(modulePath) > 0 {
		tmpDirectory := directory + ".tmp"
		pathWanted := path.Join(tmpDirectory, modulePath)

		err := os.Rename(directory, tmpDirectory)
		CheckIfError(moduleName, err)

		err = copy.Copy(pathWanted, directory)
		CheckIfError(moduleName, err)
		os.RemoveAll(tmpDirectory)
	}
	// Cleanup .git directory
	os.RemoveAll(path.Join(directory, ".git"))
}

// Handle local modules from relative paths
var localSourcePrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

func isLocalSourceAddr(addr string) bool {
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}

func copyFile(name string, src string, dst string) {
	jww.INFO.Printf("[%s] Copying from %s", name, src)
	err := copy.Copy(src, dst)
	CheckIfError(name, err)
}

// Handle modules from Terraform registy
var registryBaseURL = "https://registry.terraform.io/v1/modules"
var githubDownloadURLRe = regexp.MustCompile(`https://[^/]+/repos/([^/]+)/([^/]+)/tarball/([^/]+)/.*`)

func isRegistrySourceAddr(addr string) bool {
	nameRegex := "[0-9A-Za-z](?:[0-9A-Za-z-_]{0,62}[0-9A-Za-z])?"
	providerRegex := "[0-9a-z]{1,64}"
	registryRegex := regexp.MustCompile(
		fmt.Sprintf("^(%s)\\/(%s)\\/(%s)(?:\\/\\/(.*))?$", nameRegex, nameRegex, providerRegex))
	return registryRegex.MatchString(addr)
}

func getRegistrySource(name string, source string, version string) (string, string) {
	jww.INFO.Printf("[%s] Looking up %s version %s in Terraform registry", name, source, version)
	if version == "master" {
		err := errors.New("Registry module version must be specified")
		CheckIfError(name, err)
	}
	src := strings.Split(source, "/")
	namespace, name, provider := src[0], src[1], src[2]

	registryDownloadURL := fmt.Sprintf("%s/%s/%s/%s/%s/download",
		registryBaseURL,
		namespace,
		name,
		provider,
		version)

	client := &http.Client{}
	req, err := http.NewRequest("GET", registryDownloadURL, nil)
	CheckIfError(name, err)

	req.Header.Set("User-Agent", "XTerrafile (https://github.com/devopsmakers/xterrafile)")
	resp, err := client.Do(req)
	CheckIfError(name, err)
	defer resp.Body.Close()

	var githubDownloadURL = ""
	if len(resp.Header["X-Terraform-Get"]) > 0 {
		githubDownloadURL = resp.Header["X-Terraform-Get"][0]
	}

	if githubDownloadURLRe.MatchString(githubDownloadURL) {
		matches := githubDownloadURLRe.FindStringSubmatch(githubDownloadURL)
		user, repo, version := matches[1], matches[2], matches[3]
		source = fmt.Sprintf("github.com/%s/%s.git", user, repo)
		return source, version
	}
	err = errors.New("Unable to find module or version download url")
	CheckIfError(name, err)
	return "", "" // Never reacbhes here
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
	if IContains(source, ".git") {
		source = strings.Replace(source, "https://github.com/", "github.com/", 1)
	}

	moduleSource, err := getter.Detect(source, directory, getter.Detectors)
	CheckIfError(name, err)

	jww.DEBUG.Printf("[%s] Detected real source: %s", name, moduleSource)

	realModuleSource, err := url.Parse(moduleSource)
	CheckIfError(name, err)

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
	CheckIfError(name, err)
}

// The subDir portion will be returned as empty if no subdir separator
// ("//") is present in the address.
func splitAddrSubdir(addr string) (packageAddr, subDir string) {
	return getter.SourceDirSubdir(addr)
}
