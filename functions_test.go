package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThis(t *testing.T) {
	// Test when runtime.Caller succeeds
	funcName, fileName := this()
	assert.Equal(t, "func", funcName)
	assert.Contains(t, fileName, "TestThis")

	// Test that we get a valid function name
	assert.NotEmpty(t, funcName)
	assert.NotEmpty(t, fileName)
}
