package panic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/cff"
)

func TestCatchesPanicParallel(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsParallel()
	assert.ErrorContains(t, err, "panic: panic\nstacktrace:")
	var panicError cff.PanicError
	assert.Equal(t, errors.As(err, &panicError), true, "error returned should be a cff.PanicError")
	assert.Equal(t, panicError.Value, "panic", "PanicError.Value should be recovered value")
	assert.Contains(t, panicError.Stacktrace, "[frames]:\npanic()", "panic should show up at the top of the ")
	assert.Contains(t, panicError.Stacktrace, ".FlowPanicsParallel.func", "function that panicked should be in the stack")
}

func TestCatchesPanicSerial(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsSerial()
	assert.ErrorContains(t, err, "panic: panic\n")
	var panicError cff.PanicError
	assert.Equal(t, errors.As(err, &panicError), true, "error returned should be a cff.PanicError")
	assert.Equal(t, panicError.Value, "panic", "PanicError.Value should be recovered value")
	assert.Contains(t, panicError.Stacktrace, "[frames]:\npanic()", "panic should show up at the top of the ")
	assert.Contains(t, panicError.Stacktrace, ".FlowPanicsSerial.func", "function that panicked should be in the stack")
}
