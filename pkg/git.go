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
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

var gitSourcePrefixes = []string{
	"git::",
	"git@",
}

// IsGitSourceAddr returns true if the address is a git source
func IsGitSourceAddr(addr string) bool {
	jww.DEBUG.Printf("Testing if %s is a local source", addr)
	for _, prefix := range gitSourcePrefixes {
		if strings.HasPrefix(addr, prefix) || strings.HasSuffix(addr, ".git") {
			return true
		}
	}
	return false
}

// GetGitSource returns the source uri and version of a module from git
func GetGitSource(name string, source string, version string) (string, string) {
	var gitVersion string

	switch {
	case isConditionalVersion(version):
		var err error
		tagVersions := getGitTags(name, source)
		gitVersion, err = getModuleVersion(tagVersions, version)
		CheckIfError(name, err)
	default:
		gitVersion = version
	}
	return source, gitVersion
}

func getGitTags(name string, source string) []string {
	var stdoutbuf bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--tags", source)
	cmd.Stdout = &stdoutbuf
	err := cmd.Run()
	CheckIfError(name, err)

	var tagRegexp = regexp.MustCompile(`refs/tags/(.*)`)
	tags := []string{}

	tagScanner := bufio.NewScanner(&stdoutbuf)
	for tagScanner.Scan() {
		tag := tagRegexp.FindStringSubmatch(tagScanner.Text())
		if tag != nil {
			tags = append(tags, tag[1])
		}
	}
	return tags
}
