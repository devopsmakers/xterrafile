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
	"os"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(moduleName string, err error) {
	if err == nil {
		return
	}
	jww.FATAL.Printf("[%s] %s", moduleName, err)
	os.Exit(1)
}

// IContains is case-insensitive string contains check
func IContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}
