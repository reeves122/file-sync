package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetFiles(t *testing.T) {
	_ = os.Setenv("INPUT_FILES", "file1\nfile2\nfile3")
	assert.Equal(t, []string{"file1", "file2", "file3"}, GetFiles())
}
