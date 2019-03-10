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
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var registryBaseURL = "https://registry.terraform.io/v1/modules"
var githubDownloadURLRe = regexp.MustCompile(`https://[^/]+/repos/([^/]+)/([^/]+)/tarball/([^/]+)/.*`)

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

func getModule(moduleName string, moduleMeta module, wg *sync.WaitGroup) {
	defer wg.Done()

	moduleSource := moduleMeta.Source
	moduleVersion := "master"
	if len(moduleMeta.Version) > 0 {
		moduleVersion = moduleMeta.Version
	}
	modulePath := moduleMeta.Path

	directory := path.Join(VendorDir, moduleName)

	switch {
	case strings.HasPrefix(moduleSource, "./") || strings.HasPrefix(
		moduleSource, "../") || strings.HasPrefix(moduleSource, "/"):
		copyFile(moduleName, moduleSource, directory)
	case validRegistry(moduleSource):
		source, version := getRegistrySource(moduleName, moduleSource, moduleVersion)
		gitCheckout(moduleName, source, version, directory)
	case IContains(moduleSource, "git"):
		gitCheckout(moduleName, moduleSource, moduleVersion, directory)
	}

	// If we have a path specified, let's extract it (move and copy stuff).
	if len(modulePath) > 0 {
		tmpDirectory := directory + ".tmp"
		pathWanted := path.Join(tmpDirectory, modulePath)

		err := os.Rename(directory, tmpDirectory)
		CheckIfError(err)

		err = copy.Copy(pathWanted, directory)
		CheckIfError(err)
		os.RemoveAll(tmpDirectory)
	}
	// Cleanup .git directoriy
	os.RemoveAll(path.Join(directory, ".git"))
}

func init() {
	jww.SetStdoutThreshold(jww.LevelInfo)
	rootCmd.AddCommand(installCmd)
}

func getRegistrySource(name string, source string, version string) (string, string) {
	jww.INFO.Printf("[%s] Looking up %s version %s in Terraform registry", name, source, version)
	if version == "master" {
		err := errors.New("Registry module version must be specified")
		CheckIfError(err)
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
	CheckIfError(err)

	req.Header.Set("User-Agent", "XTerrafile (https://github.com/devopsmakers/xterrafile)")
	resp, err := client.Do(req)
	CheckIfError(err)
	defer resp.Body.Close()

	var githubDownloadURL = ""
	if len(resp.Header["X-Terraform-Get"]) > 0 {
		githubDownloadURL = resp.Header["X-Terraform-Get"][0]
	}

	if githubDownloadURLRe.MatchString(githubDownloadURL) {
		matches := githubDownloadURLRe.FindStringSubmatch(githubDownloadURL)
		user, repo, version := matches[1], matches[2], matches[3]
		source = fmt.Sprintf("https://github.com/%s/%s.git", user, repo)
		return source, version
	}
	err = errors.New("Unable to find module / version download url")
	CheckIfError(err)
	return "", "" // Never reacbhes here
}

func validRegistry(source string) bool {
	nameRegex := "[0-9A-Za-z](?:[0-9A-Za-z-_]{0,62}[0-9A-Za-z])?"
	providerRegex := "[0-9a-z]{1,64}"
	registryRegex := regexp.MustCompile(
		fmt.Sprintf("^(%s)\\/(%s)\\/(%s)(?:\\/\\/(.*))?$", nameRegex, nameRegex, providerRegex))
	return registryRegex.MatchString(source)
}

func copyFile(name string, src string, dst string) {
	jww.INFO.Printf("[%s] Copying from %s", name, src)
	err := copy.Copy(src, dst)
	CheckIfError(err)
}

func gitCheckout(name string, repo string, version string, directory string) {
	jww.INFO.Printf("[%s] Checking out %s from %s", name, version, repo)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: repo,
	})
	CheckIfError(err)

	h, err := r.ResolveRevision(plumbing.Revision(version))
	if err != nil {
		h, err = r.ResolveRevision(plumbing.Revision("origin/" + version))
	}
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(h.String()),
	})
	CheckIfError(err)
}
