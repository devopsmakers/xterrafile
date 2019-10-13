package xterrafile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModuleVersion(t *testing.T) {
	var version string
	var err error

	version, _ = getModuleVersion([]string{"1.1.1", "2.1.1", "2.0.1"}, "1.1.1")
	assert.Equal(t, "1.1.1", version, "version should be >= 2.0.0 < 2.2.0")

	version, _ = getModuleVersion([]string{"1.1.1", "not-a-version", "2.1.1", "2.0.1"}, ">= 2.0.0 < 2.2.0")
	assert.Equal(t, "2.1.1", version, "version should be 1.1.1")

	_, err = getModuleVersion([]string{"1.1.1", "2.1.1", "2.0.1"}, ">= no < version")
	assert.EqualError(t, err, "Could not get version from string: \">=no\"")

	_, err = getModuleVersion([]string{"not", "a", "version"}, ">= 2.0.0 < 2.2.0")
	assert.EqualError(t, err, "Unable to find a valid version of this module")
}
