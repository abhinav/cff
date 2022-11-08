package panic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatchesPanicParallel(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsParallel()
	assert.ErrorContains(t, err, "panic: panic\nstacktrace:")
}

func TestCatchesPanicSerial(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsSerial()
	assert.ErrorContains(t, err, "panic: panic\nstacktrace:")
}
