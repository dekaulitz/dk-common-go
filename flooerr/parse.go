package flooerr

import (
	"core-common-go/flooerr/internal"
	"errors"
)

// ErrorInfo contains all information extracted from a FlooErr
type ErrorInfo struct {
	Code       internal.Code
	Message    string
	ErrorMsg   string
	Context    map[string]any
	SDC        map[string]string
	StackTrace []stacktrace
	Cause      error
	IsFlooErr  bool
}

// Parse extracts all information from an error.
// If the error is a FlooErr, it returns detailed information.
// Otherwise, it returns basic information with IsFlooErr = false.
func Parse(err error) ErrorInfo {
	if err == nil {
		return ErrorInfo{
			IsFlooErr: false,
		}
	}

	flooErr, ok := AsFlooErr(err)
	if !ok {
		return ErrorInfo{
			ErrorMsg:  err.Error(),
			IsFlooErr: false,
		}
	}

	return ErrorInfo{
		Code:       flooErr.Code(),
		Message:    flooErr.Message(),
		ErrorMsg:   flooErr.Error(),
		Context:    flooErr.Context(),
		SDC:        flooErr.SDC(),
		StackTrace: flooErr.StackTrace(),
		Cause:      flooErr.Unwrap(),
		IsFlooErr:  true,
	}
}

// AsFlooErr checks if an error is a FlooErr and returns it.
// Returns false if the error is nil or not a FlooErr.
func AsFlooErr(err error) (FlooErr, bool) {
	if err == nil {
		return nil, false
	}

	var flooErr FlooErr
	if errors.As(err, &flooErr) {
		return flooErr, true
	}

	return nil, false
}

// GetCode extracts the error code from an error.
// Returns empty string if the error is not a FlooErr.
func GetCode(err error) internal.Code {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return internal.Code("")
	}
	return flooErr.Code()
}

// GetCodeString extracts the error code as a string from an error.
// Returns empty string if the error is not a FlooErr.
func GetCodeString(err error) string {
	return GetCode(err).String()
}

// GetMessage extracts the error message from an error.
// Returns empty string if the error is not a FlooErr.
func GetMessage(err error) string {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return ""
	}
	return flooErr.Message()
}

// GetContext extracts the context map from an error.
// Returns nil if the error is not a FlooErr.
func GetContext(err error) map[string]any {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return nil
	}
	return flooErr.Context()
}

// GetContextValue extracts a specific context value by key from an error.
// Returns nil if the error is not a FlooErr or the key doesn't exist.
func GetContextValue(err error, key string) any {
	context := GetContext(err)
	if context == nil {
		return nil
	}
	return context[key]
}

// GetSDC extracts the SDC (Structured Diagnostic Context) map from an error.
// Returns nil if the error is not a FlooErr.
func GetSDC(err error) map[string]string {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return nil
	}
	return flooErr.SDC()
}

// GetSDCValue extracts a specific SDC value by key from an error.
// Returns empty string if the error is not a FlooErr or the key doesn't exist.
func GetSDCValue(err error, key string) string {
	sdc := GetSDC(err)
	if sdc == nil {
		return ""
	}
	return sdc[key]
}

// GetStackTrace extracts the stack trace from an error.
// Returns nil if the error is not a FlooErr or stack trace is not available.
func GetStackTrace(err error) []stacktrace {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return nil
	}
	return flooErr.StackTrace()
}

// GetCause extracts the underlying cause error from an error.
// Returns nil if the error is not a FlooErr or has no cause.
func GetCause(err error) error {
	flooErr, ok := AsFlooErr(err)
	if !ok {
		return nil
	}
	return flooErr.Unwrap()
}

// IsFlooErr checks if an error is a FlooErr.
func IsFlooErr(err error) bool {
	_, ok := AsFlooErr(err)
	return ok
}

// UnwrapChain unwraps the entire error chain and returns all errors in the chain.
// The first element is the top-level error, and the last is the root cause.
func UnwrapChain(err error) []error {
	if err == nil {
		return nil
	}

	var chain []error
	current := err

	for current != nil {
		chain = append(chain, current)

		flooErr, ok := AsFlooErr(current)
		if ok {
			current = flooErr.Unwrap()
		} else {
			// Try standard errors.Unwrap for non-FlooErr errors
			current = errors.Unwrap(current)
		}
	}

	return chain
}

// GetRootCause returns the root cause error (the deepest error in the chain).
func GetRootCause(err error) error {
	chain := UnwrapChain(err)
	if len(chain) == 0 {
		return nil
	}
	return chain[len(chain)-1]
}

// HasCode checks if an error has a specific error code.
func HasCode(err error, code string) bool {
	return GetCodeString(err) == code
}

// HasContextKey checks if an error has a specific context key.
func HasContextKey(err error, key string) bool {
	context := GetContext(err)
	if context == nil {
		return false
	}
	_, exists := context[key]
	return exists
}

// HasSDCKey checks if an error has a specific SDC key.
func HasSDCKey(err error, key string) bool {
	sdc := GetSDC(err)
	if sdc == nil {
		return false
	}
	_, exists := sdc[key]
	return exists
}
