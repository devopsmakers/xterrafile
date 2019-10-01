package xterrafile

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLocalSourceAddr(t *testing.T) {
	assert.True(t, IsLocalSourceAddr("./some/path"), "relative to current directory should be true")
	assert.True(t, IsLocalSourceAddr("../some/path"), "relative to current directory should be true")
	assert.False(t, IsLocalSourceAddr("/some/absolute/path"), "absolute path should be false")
	assert.False(t, IsLocalSourceAddr("http://something"), "http source should be false")
}

func TestCopyFile(t *testing.T) {
	tmpTestDir := "../../test/tmp/"
	moduleName := "test-module"
	os.RemoveAll(tmpTestDir)
	CopyFile("test-dir", "../../test/module", path.Join(tmpTestDir, moduleName))
	assert.FileExists(t, path.Join(tmpTestDir, moduleName, "main.tf"), "file should be copied into tmp dir")

}
