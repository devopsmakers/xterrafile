package xterrafile

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jww "github.com/spf13/jwalterweatherman"
)

func TestCheckIfError(t *testing.T) {
	// Capture logging
	var outputBuf bytes.Buffer
	jww.SetStdoutOutput(&outputBuf)
	defer jww.SetStdoutOutput(os.Stdout)

	// override osExit to test for usage
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()
	myExit := func(code int) {
		panic("osExit called")
	}
	osExit = myExit

	assert.PanicsWithValue(t, "osExit called",
		func() { CheckIfError("test-error", errors.New("A test Error")) }, "os.Exit should be called")

	require.Contains(t, outputBuf.String(), "FATAL")
	require.Contains(t, outputBuf.String(), "test-error")
	require.Contains(t, outputBuf.String(), "A test Error")

	assert.NotPanics(t,
		func() { CheckIfError("test-no-error", nil) }, "os.Exit should not be called")
}

func TestIContains(t *testing.T) {
	assert.True(t, IContains("teststring", "TEST"), "string comparison should be true")
	assert.False(t, IContains("TEST", "teststring"), "string comparison should be false")
}
