package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSerial(t *testing.T) {
	TestSerial := NewSerial("COM5", 230400, 8, 1, time.Second*5)

	ConfMade := true  // Should be a valid config
	PortOpen := false // Port should be false as it hasn't been called

	assert.Equal(t, TestSerial.ConfigMade, ConfMade, "ConfigMade should be true")
	assert.Equal(t, TestSerial.PortOpen, PortOpen, "PortOpen should be false")
}
