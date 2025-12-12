package internal

import (
	"fmt"
	"runtime"
)

type Code string

func (code Code) String() string {
	return string(code)
}

type ErrProps struct {
	message        string
	code           string
	withStackTrace bool
	context        map[string]any
	sdc            map[string]string
}

func create() *ErrProps {
	context := make(map[string]any)
	sdc := make(map[string]string)
	return &ErrProps{
		context:        context,
		sdc:            sdc,
		withStackTrace: true,
	}
}

// Create creates a new ErrProps instance
func Create() *ErrProps {
	return create()
}

// WithMessage sets the message for the error
func (receiver *ErrProps) WithMessage(message string) *ErrProps {
	receiver.message = message
	return receiver
}

func (receiver *ErrProps) WithCode(code string) *ErrProps {
	receiver.code = code
	return receiver
}

func (receiver *ErrProps) WithStackTrace(enableStackTrace bool) *ErrProps {
	receiver.withStackTrace = enableStackTrace
	return receiver
}

func (receiver *ErrProps) WithContext(key string, value any) *ErrProps {
	receiver.context[key] = value
	return receiver
}

func (receiver *ErrProps) WithSDC(key string, value string) *ErrProps {
	receiver.sdc[key] = value
	return receiver
}

func (receiver *ErrProps) Build(cause error, message string) error {
	var stackTracePTR []uintptr
	if receiver.withStackTrace {
		stackTracePTR = callers(4) // Skip Build, caller of Build, and runtime frames
	}

	errMessage := message
	if receiver.message != "" {
		errMessage = receiver.message
	}

	// Create error using BuildErr function which should be set by flooerr package
	if buildErrFunc != nil {
		return buildErrFunc(
			receiver.message,
			errMessage,
			Code(receiver.code),
			cause,
			stackTracePTR,
			receiver.context,
			receiver.sdc,
		)
	}

	// Fallback: return a simple error if builder is not set
	if cause != nil {
		return cause
	}
	return &simpleError{message: errMessage}
}

// Error creates and returns an error with the configured properties.
// If cause is provided, it will be wrapped as the underlying error.
// The message parameter is used as a fallback if no message was set via WithMessage().
func (receiver *ErrProps) Error(cause error, message string) error {
	return receiver.Build(cause, message)
}

func (receiver *ErrProps) Errorf(format string, args ...any) error {
	return receiver.Build(nil, fmt.Sprintf(format, args...))
}

func (receiver *ErrProps) Wrap(cause error, message string) error {
	return receiver.Build(cause, message)
}

func (receiver *ErrProps) Wrapf(cause error, format string, args ...any) error {
	return receiver.Build(cause, fmt.Sprintf(format, args...))
}

// BuildErrFunc is a function type for building errors from internal package
type BuildErrFunc func(
	message string,
	errMessage string,
	code Code,
	cause error,
	stackTracePTR []uintptr,
	context map[string]any,
	sdc map[string]string,
) error

var buildErrFunc BuildErrFunc

// SetBuildErrFunc sets the error builder function (called from flooerr package)
func SetBuildErrFunc(fn BuildErrFunc) {
	buildErrFunc = fn
}

// simpleError is a fallback error implementation
type simpleError struct {
	message string
}

func (e *simpleError) Error() string {
	return e.message
}

func callers(skip int) []uintptr {
	const depth = 15
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st = pcs[0 : n-2]
	return st
}
