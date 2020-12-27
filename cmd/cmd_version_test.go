package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	var tests = []struct {
		cmd string
	}{
		{"version"},
		{"ver"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		buf := bytes.NewBufferString("")
		appWriter = buf

		// Act
		_ = Execute(test.cmd)

		// Assert
		ass.Contains(buf.String(), Version)
	}
}
