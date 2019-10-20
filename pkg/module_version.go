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
	"errors"
	"sort"

	"github.com/blang/semver"
)

func isValidVersion(version string) bool {
	_, err := semver.ParseTolerant(version)
	if err != nil {
		return false
	}
	return true
}

func isConditionalVersion(versionConditional string) bool {
	_, err := semver.ParseRange(versionConditional)
	if err != nil {
		return false
	}
	return true
}

func getModuleVersion(sourceVersions []string, versionConditional string) (string, error) {
	var validSourceVersions []semver.Version
	var originalVersions []string

	for _, sourceVersion := range sourceVersions {
		v, err := semver.ParseTolerant(sourceVersion)
		if err != nil {
			continue
		}
		validSourceVersions = append(validSourceVersions, v)
		originalVersions = append(originalVersions, sourceVersion)
	}

	semver.Sort(validSourceVersions)
	sort.Strings(originalVersions)

	// Assume latest if we get passed an empty string
	if versionConditional == "" {
		return originalVersions[len(originalVersions)-1], nil
	}

	validModuleVersionRange, err := semver.ParseRange(versionConditional)
	if err != nil {
		return "", err
	}

	for i := range validSourceVersions {
		v := validSourceVersions[len(validSourceVersions)-1-i]
		o := originalVersions[len(originalVersions)-1-i]
		if validModuleVersionRange(v) {
			return o, nil
		}
	}

	err = errors.New("Unable to find a valid version of this module")
	return "", err
}
