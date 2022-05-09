package common

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_RemoveDir(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test")
	assert.NotPanics(t, func() {
		RemoveDir(dir)
	})
}

func Test_RunCommand_Success(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test")
	output, err := RunCommand(dir, "echo", "foo")
	assert.Contains(t, output, "foo")
	assert.NoError(t, err)
}

func Test_RunCommand_Error(t *testing.T) {
	dir, _ := ioutil.TempDir("", "test")
	output, err := RunCommand(dir, "foo", "foo")
	assert.Contains(t, output, "")
	assert.Error(t, err)
}

func Test_LogCommand(t *testing.T) {
	assert.NotPanics(t, func() {
		LogCommand("cd", "foo")
	})
}

func Test_LogOutput(t *testing.T) {
	var test bytes.Buffer
	test.Write([]byte("test\n"))

	assert.NotPanics(t, func() {
		LogOutput(test)
	})
}

func Test_LogOutput_Empty(t *testing.T) {
	var test bytes.Buffer
	test.Write([]byte(""))

	assert.NotPanics(t, func() {
		LogOutput(test)
	})
}
