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
		"Checking out master from git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v1.46.0 from git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out 01601169c00c68f37d5df8a80cc17c88f02c04d0 from git@github.com:terraform-aws-modules/terraform-aws-vpc",
		"Checking out v0.7.0 from https://github.com/claranet/terraform-aws-lambda.git",
		"Copying from ./test/module",
		"Looking up claranet/lambda/aws version 0.7.0 in Terraform registry",
	} {
		assert.Contains(t, testcli.Stdout(), output)
	}
	// Assert files exist
	for _, moduleName := range []string{
		"tf-aws-vpc",
		"tf-aws-vpc-experimental",
		"tf-aws-vpc-commit",
		"tf-aws-vpc-default",
		"terraform-aws-lambda",
		"terraform-test-path",
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
	var yaml = `tf-aws-vpc:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "v1.46.0"
tf-aws-vpc-experimental:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "master"
tf-aws-vpc-commit:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: "01601169c00c68f37d5df8a80cc17c88f02c04d0"
tf-aws-vpc-default:
  source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
terraform-aws-lambda:
  source: "claranet/lambda/aws"
  version: "0.7.0"
terraform-test-path:
  source: "./test/module"
`
	createFile(t, path.Join(folder, "Terrafile.test"), yaml)
}
