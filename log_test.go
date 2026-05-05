package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFieldName(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string field", "test", "test="},
		{"number field", 123, "%!s(int=123)="},
		{"empty field", "", "="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFieldName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigureLogging(t *testing.T) {
	// Setup test config
	G = &Global{
		Config: &Config{},
	}
	G.Config.Log.Format = "text"
	G.Config.Log.Color = false

	// Test with valid log level
	logger, err := configureLogging("test", "info")
	assert.NoError(t, err)
	assert.NotZero(t, logger)

	// Test with invalid log level
	_, err = configureLogging("test", "invalid")
	assert.Error(t, err)

	// Test with production environment
	G.Config.Log.Format = "json"
	logger, err = configureLogging("prod", "info")
	assert.NoError(t, err)
	assert.NotZero(t, logger)
}
