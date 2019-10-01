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
	"strings"

	"github.com/otiai10/copy"
	jww "github.com/spf13/jwalterweatherman"
)

// Handle local modules from relative paths
var localSourcePrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

// IsLocalSourceAddr returns true if the address is a relative filepath
func IsLocalSourceAddr(addr string) bool {
	for _, prefix := range localSourcePrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}
	return false
}

// CopyFile copies files to the correct destination
func CopyFile(name string, src string, dst string) {
	jww.INFO.Printf("[%s] Copying from %s", name, src)
	err := copy.Copy(src, dst)
	CheckIfError(name, err)
}
