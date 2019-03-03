package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

var xTerrafileBinaryPath string
var workingDirectory string

func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	xTerrafileBinaryPath = workingDirectory + "/xterrafile"
}
func TestTerraformWithTerrafilePath(t *testing.T) {
	folder, back := setup(t)
	defer back()

	testcli.Run(xTerrafileBinaryPath, "-f", fmt.Sprint(folder, "/Terrafile.test"), "install")

	if !testcli.Success() {
		t.Fatalf("Expected to succeed, but failed: %q with message: %q", testcli.Error(), testcli.Stderr())
	}
	// Assert output
	for _, output := range []string{
		"Removing all modules in vendor/xterrafile",
		"[terrafile-test-registry] Looking up terraform-digitalocean-modules/droplet/digitalocean version 0.1.7 in Terraform registry",
		"[terrafile-test-registry] Checking out v0.1.7 from https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-https] Checking out master from https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-tag] Checking out v0.1.7 from git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-branch] Checking out branch_test from git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-commit] Checking out 2e6b9729f3f6ea3ef5190bac0b0e1544a01fd80f from https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-path] Checking out v0.1.7 from https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git",
		"[terrafile-test-local] Copying from ./test/module",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert files exist
	for _, moduleName := range []string{
		"terrafile-test-registry",
		"terrafile-test-https",
		"terrafile-test-tag",
		"terrafile-test-branch",
		"terrafile-test-commit",
		"terrafile-test-path",
		"terrafile-test-local",
	} {
		assert.DirExists(t, path.Join(workingDirectory, "vendor/xterrafile", moduleName))
	}
}

func setup(t *testing.T) (current string, back func()) {
	folder, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	createTerrafile(t, folder)
	return folder, func() {
		assert.NoError(t, os.RemoveAll(folder))
	}
}

func createFile(t *testing.T, filename string, contents string) {
	assert.NoError(t, ioutil.WriteFile(filename, []byte(contents), 0644))
}

func createTerrafile(t *testing.T, folder string) {
	var yaml = `terrafile-test-registry:
  source: "terraform-digitalocean-modules/droplet/digitalocean"
  version: "0.1.7"
terrafile-test-https:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
terrafile-test-tag:
  source: "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "v0.1.7"
terrafile-test-branch:
  source: "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "branch_test"
terrafile-test-commit:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "2e6b9729f3f6ea3ef5190bac0b0e1544a01fd80f"
terrafile-test-path:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "v0.1.7"
  path: "examples/simple"
terrafile-test-local:
  source: "./test/module"
`
	createFile(t, path.Join(folder, "Terrafile.test"), yaml)
}
