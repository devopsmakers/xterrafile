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
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"

	jww "github.com/spf13/jwalterweatherman"
)

// IsRegistrySourceAddr check an address is a valid registry address
func IsRegistrySourceAddr(addr string) bool {
	jww.DEBUG.Printf("Testing if %s is a registry source", addr)
	_, err := regsrc.ParseModuleSource(addr)
	return err == nil
}

// GetRegistrySource retrieves a modules download source from a Terraform registry
func GetRegistrySource(name string, source string, version string, services *disco.Disco) (string, string) {
	var modVersions []string

	modSrc, err := getModSrc(source)
	CheckIfError(name, err)

	regClient := registry.NewClient(services, nil)

	switch {
	case isValidVersion(version):
		_ = version
	default:
		modVersions = getRegistryModuleVersions(regClient, modSrc)
		version, err = getModuleVersion(modVersions, version)
		CheckIfError(name, err)
	}

	jww.INFO.Printf("[%s] Found module version %s at %s", name, version, modSrc.Host())

	regSrc, err := regClient.ModuleLocation(modSrc, version)
	CheckIfError(name, err)
	jww.INFO.Printf("[%s] Downloading from source URL %s", name, regSrc)

	return regSrc, version
}

// Helper function to return a list of available module versions
func getRegistryModuleVersions(c *registry.Client, modSrc *regsrc.Module) []string {
	// Don't log from Terraform's HTTP client
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	regClientResp, _ := c.ModuleVersions(modSrc)

	regModule := regClientResp.Modules[0]
	moduleVersions := []string{}

	for _, moduleVersion := range regModule.Versions {
		moduleVersions = append(moduleVersions, moduleVersion.Version)
	}

	return moduleVersions
}

// Helper function to parse and return a module source
func getModSrc(source string) (*regsrc.Module, error) {
	modSrc, err := regsrc.ParseModuleSource(source)
	if err != nil {
		return nil, err
	}
	return modSrc, nil
}
