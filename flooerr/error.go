package flooerr

import (
	"core-common-go/flooerr/internal"
	"fmt"
	"runtime"
)

type FlooErr interface {
	error
	Code() internal.Code
	Message() string
	StackTrace() []stacktrace
	Unwrap() error
	Context() map[string]any
	SDC() map[string]string
}

type err struct {
	message       string
	errMessage    string
	code          internal.Code
	cause         error
	stackTracePTR []uintptr
	stackTrace    []stacktrace
	context       map[string]any
	sdc           map[string]string
}

func (e *err) Code() internal.Code {
	return e.code
}

func (e *err) Message() string {
	return e.message
}

func (e *err) StackTrace() []stacktrace {
	if len(e.stackTrace) > 0 {
		return e.stackTrace
	}
	if len(e.stackTracePTR) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(e.stackTracePTR)
	var traces []stacktrace

	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		traces = append(traces, stacktrace{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
	}

	e.stackTrace = traces
	return traces
}

func (e *err) Unwrap() error {
	return e.cause
}

func (e *err) Context() map[string]any {
	return e.context
}

func (e *err) SDC() map[string]string {
	return e.sdc
}

func (e *err) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s; caused by: %v", e.errMessage, e.cause)
	}
	return e.errMessage
}

func Message(msg string) *internal.ErrProps {
	return internal.Create().WithMessage(msg)
}

func Code(code internal.Code) *internal.ErrProps {
	return internal.Create().WithCode(code.String())
}

func StackTrace() *internal.ErrProps {
	return internal.Create().WithStackTrace(true)
}

func Context(key string, value any) *internal.ErrProps {
	return internal.Create().WithContext(key, value)
}

func SDC(key string, value map[string]string) *internal.ErrProps {
	return internal.Create().WithContext(key, value)
}

func Error(message string) error {
	return internal.Create().Build(nil, message)
}

func ErrorF(format string, args ...any) error {
	return internal.Create().Build(nil, fmt.Sprintf(format, args...))
}

func Wrap(err error, message string) error {
	return internal.Create().Build(err, message)
}

func WrapF(err error, format string, args ...any) error {
	return internal.Create().Build(err, fmt.Sprintf(format, args...))
}

// newErr creates a new err struct. This is used by internal package.
func newErr(
	message string,
	errMessage string,
	code internal.Code,
	cause error,
	stackTracePTR []uintptr,
	context map[string]any,
	sdc map[string]string,
) FlooErr {
	return &err{
		message:       message,
		errMessage:    errMessage,
		code:          code,
		cause:         cause,
		stackTracePTR: stackTracePTR,
		stackTrace:    nil,
		context:       context,
		sdc:           sdc,
	}
}

func init() {
	internal.SetBuildErrFunc(func(
		message string,
		errMessage string,
		code internal.Code,
		cause error,
		stackTracePTR []uintptr,
		context map[string]any,
		sdc map[string]string,
	) error {
		return newErr(message, errMessage, code, cause, stackTracePTR, context, sdc)
	})
}
