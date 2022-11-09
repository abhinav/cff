// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cff

import (
	"fmt"
	"runtime"
)

// PanicError is a custom error that is thrown when a task panics. It contains the value
// that is recovered from the panic and the stacktrace of where the panic happened.
// For example, with the following flow:
//
//	var r string
//	err := cff.Flow(
//		context.Background(),
//		cff.Results(&r),
//		cff.Task(
//			func() string {
//				panic("panic")
//			},
//		),
//	)
//
// You can do the following to determine if an error from the cff.Flow is from a panic:
//
//	var panicError cff.PanicError
//	if errors.As(err, &panicError) {
//		// err is from a panic
//		fmt.Printf("recovered: %s\n", panicError.Value)
//	} else {
//		// err is not from a panic
//	}
//

type PanicError struct {
	Value      any
	Stacktrace string
}

var _ error = PanicError{}

func (pe PanicError) Error() string {
	return fmt.Sprintf("panic: %v\nstacktrace:\n%s", pe.Value, pe.Stacktrace)
}

// NewPanicError intializes a PanicError with the given value and a stack trace
// that starts from the call to panic.
func NewPanicError(value any) PanicError {
	return PanicError{
		Value:      value,
		Stacktrace: panicStacktrace(),
	}
}

// panicStacktrace traverses a list of callers in the stack and finds where panic
// happened and returns a stacktrace string starting from panic
func panicStacktrace() string {
	pc := make([]uintptr, 20)
	// skipping 3 in the callers, which are:
	// - Callers
	// - caller of Callers (self)
	// - caller of self (NewPanicError)
	// becase panic should not be in any of these callers
	n := runtime.Callers(3, pc)
	if n == 0 {
		return ""
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	seenPanic := false
	stacktrace := "[frames]:\n"

	for {
		frame, more := frames.Next()
		if frame.Function == "runtime.gopanic" {
			seenPanic = true
		}
		if seenPanic {
			stacktrace = fmt.Sprintf("%s%s\n", stacktrace, formatFrame(frame))
		}
		// Check whether there are more frames to process after this one
		// or if stack trace is getting too long. 1024 was chosen to match
		// length used by runtime/debug.Stack()
		if !more || len(stacktrace) >= 1024 {
			break
		}
	}

	return stacktrace
}

func formatFrame(frame runtime.Frame) string {
	funcName := frame.Function
	if funcName == "runtime.gopanic" {
		funcName = "panic"
	}
	return fmt.Sprintf("%s()\n\t%s:%d", funcName, frame.File, frame.Line)
}
