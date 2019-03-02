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

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the modules in your Terrafile",
	Run: func(cmd *cobra.Command, args []string) {
		os.RemoveAll(VendorDir)
		os.MkdirAll(VendorDir, os.ModePerm)

		for moduleName, moduleMeta := range Config {
			moduleSource := moduleMeta.Source
			moduleVersion := "master"
			if len(moduleMeta.Version) > 0 {
				moduleVersion = moduleMeta.Version
			}

			switch {
			case IContains(moduleSource, "git"):
				gitCheckout(moduleName, moduleSource, moduleVersion)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func gitCheckout(name string, repo string, version string) {
	jww.WARN.Printf("[%s] Checking out %s from %s", name, version, repo)

	directory := path.Join(VendorDir, name)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:        repo,
		NoCheckout: true,
	})
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// Try checkoing out commits, tags and branches
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(version),
	})
	if err != nil {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/tags/" + version),
		})
	}
	if err != nil {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/heads/" + version),
		})
	}
	CheckIfError(err)

	// Cleanup .git directory
	os.RemoveAll(path.Join(directory, ".git"))
}
