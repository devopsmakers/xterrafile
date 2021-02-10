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

type semverMap struct {
	semver semver.Version
	original string
}

type semverMapList []semverMap

func (s semverMapList) Len() int {
	return len(s)
}

func (s semverMapList) Less(i, j int) bool {
	return s[i].semver.GT(s[j].semver)
}

func (s semverMapList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func isValidVersion(version string) bool {
	_, err := semver.ParseTolerant(version)

	return err == nil
}

func isConditionalVersion(versionConditional string) bool {
	_, err := semver.ParseRange(versionConditional)

	return err == nil
}

func getModuleVersion(sourceVersions []string, versionConditional string) (string, error) {
	var versions = make(semverMapList, 0, len(sourceVersions))

	for _, sourceVersion := range sourceVersions {
		v, err := semver.ParseTolerant(sourceVersion)
		if err != nil {
			// todo log something
			continue
		}
		versions = append(versions, semverMap{semver: v, original: sourceVersion})
	}

	sort.Sort(versions)

	if len(versions) == 0 {
		return "", errors.New("unable to find a valid version of this module")
	}

	// Assume latest if we get passed an empty string
	if versionConditional == "" {
		return versions[0].original, nil
	}

	validModuleVersionRange, err := semver.ParseRange(versionConditional)
	if err != nil {
		return "", err
	}

	for _, version := range versions {
		if validModuleVersionRange(version.semver) {
			return version.original, nil
		}
	}

	return "", errors.New("unable to find a valid version of this module")
}
