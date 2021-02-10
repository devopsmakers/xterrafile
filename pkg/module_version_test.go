package xterrafile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidVersion(t *testing.T) {
	assert.False(t, isValidVersion("2e6b9729f3f6ea3ef5190bac0b0e1544a01fd80f"))
	assert.False(t, isValidVersion(">= 2.0.0 < 2.2.0"))
	assert.True(t, isValidVersion("1.1.1"))
}

func TestIsConditionalVersion(t *testing.T) {
	assert.False(t, isConditionalVersion("2e6b9729f3f6ea3ef5190bac0b0e1544a01fd80f"))
	assert.True(t, isConditionalVersion(">= 2.0.0 < 2.2.0"))
	assert.True(t, isConditionalVersion("1.1.1"))
}

func TestGetModuleVersion(t *testing.T) {
	var version string
	var err error

	version, _ = getModuleVersion([]string{"v2.9.0","v2.10.0","v2.65.0","v2.66.0"}, "2.66.0") // https://github.com/devopsmakers/xterrafile/issues/30
	assert.Equal(t, "v2.66.0", version, "version should be v2.66.0")

	version, _ = getModuleVersion([]string{"1.1.1", "2.1.1", "2.0.1"}, "1.1.1")
	assert.Equal(t, "1.1.1", version, "version should be 1.1.1")

	version, _ = getModuleVersion([]string{"1.1.1", "2.1.1", "2.0.1"}, "")
	assert.Equal(t, "2.1.1", version, "version should be 2.1.1")

	version, _ = getModuleVersion([]string{"1.1.1", "not-a-version", "2.1.1", "2.0.1"}, ">= 2.0.0 < 2.2.0")
	assert.Equal(t, "2.1.1", version, "version should be 2.1.1")

	_, err = getModuleVersion([]string{"1.1.1", "2.1.1", "2.0.1"}, ">= no < version")
	assert.EqualError(t, err, "Could not get version from string: \">=no\"")

	_, err = getModuleVersion([]string{"not", "a", "version"}, ">= 2.0.0 < 2.2.0")
	assert.EqualError(t, err, "unable to find a valid version of this module")

	_, err = getModuleVersion([]string{}, ">= 2.0.0 < 2.2.0")
	assert.EqualError(t, err, "unable to find a valid version of this module")

	_, err = getModuleVersion([]string{}, "")
	assert.EqualError(t, err, "unable to find a valid version of this module")
}
