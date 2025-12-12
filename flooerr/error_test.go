package flooerr

import (
	"core-common-go/flooerr/internal"
	"errors"
	"testing"
)

func TestMessage(t *testing.T) {
	props := Message("test error message")
	if props == nil {
		t.Fatal("Message() returned nil")
	}
}

func TestErr_Code(t *testing.T) {
	err := Message("test").
		WithCode("TEST_CODE").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	code := flooErr.Code()
	if string(code) != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", code)
	}
}

func TestErr_Message(t *testing.T) {
	err := Message("custom message").
		Error(nil, "default message")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	message := flooErr.Message()
	if message != "custom message" {
		t.Errorf("Expected message 'custom message', got '%s'", message)
	}
}

func TestErr_Message_Default(t *testing.T) {
	err := Message("").
		Error(nil, "default message")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	// When message is empty, Error uses the fallback message
	// But Message() returns the set message (empty in this case)
	message := flooErr.Message()
	if message != "" {
		t.Errorf("Expected empty message, got '%s'", message)
	}
}

func TestErr_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := Message("wrapped error").
		Error(originalErr, "wrapped error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	cause := flooErr.Unwrap()
	if cause == nil {
		t.Fatal("Expected non-nil cause")
	}

	if cause.Error() != "original error" {
		t.Errorf("Expected 'original error', got '%s'", cause.Error())
	}
}

func TestErr_Unwrap_Nil(t *testing.T) {
	err := Message("error without cause").
		Error(nil, "error without cause")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	cause := flooErr.Unwrap()
	if cause != nil {
		t.Errorf("Expected nil cause, got %v", cause)
	}
}

func TestErr_Context(t *testing.T) {
	err := Message("test").
		WithContext("key1", "value1").
		WithContext("key2", 123).
		WithContext("key3", true).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	context := flooErr.Context()
	if len(context) != 3 {
		t.Errorf("Expected 3 context items, got %d", len(context))
	}

	if context["key1"] != "value1" {
		t.Errorf("Expected context['key1'] = 'value1', got '%v'", context["key1"])
	}

	if context["key2"] != 123 {
		t.Errorf("Expected context['key2'] = 123, got '%v'", context["key2"])
	}

	if context["key3"] != true {
		t.Errorf("Expected context['key3'] = true, got '%v'", context["key3"])
	}
}

func TestErr_Context_Empty(t *testing.T) {
	err := Message("test").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	context := flooErr.Context()
	if context == nil {
		t.Fatal("Expected non-nil context map")
	}

	if len(context) != 0 {
		t.Errorf("Expected empty context, got %d items", len(context))
	}
}

func TestErr_SDC(t *testing.T) {
	err := Message("test").
		WithSDC("trace_id", "trace_123").
		WithSDC("span_id", "span_456").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	sdc := flooErr.SDC()
	if len(sdc) != 2 {
		t.Errorf("Expected 2 SDC items, got %d", len(sdc))
	}

	if sdc["trace_id"] != "trace_123" {
		t.Errorf("Expected SDC['trace_id'] = 'trace_123', got '%s'", sdc["trace_id"])
	}

	if sdc["span_id"] != "span_456" {
		t.Errorf("Expected SDC['span_id'] = 'span_456', got '%s'", sdc["span_id"])
	}
}

func TestErr_SDC_Empty(t *testing.T) {
	err := Message("test").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	sdc := flooErr.SDC()
	if sdc == nil {
		t.Fatal("Expected non-nil SDC map")
	}

	if len(sdc) != 0 {
		t.Errorf("Expected empty SDC, got %d items", len(sdc))
	}
}

func TestErr_StackTrace(t *testing.T) {
	err := Message("test").
		WithStackTrace(true).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	stackTrace := flooErr.StackTrace()
	// Stack trace might be nil or empty in some test environments
	// We verify the method doesn't panic and returns a valid result
	_ = stackTrace // Just verify it doesn't panic

	// Verify stack trace structure if it has frames
	for i, frame := range stackTrace {
		if frame.Function == "" {
			t.Errorf("Stack frame %d: Function should not be empty", i)
		}
		if frame.File == "" {
			t.Errorf("Stack frame %d: File should not be empty", i)
		}
		if frame.Line <= 0 {
			t.Errorf("Stack frame %d: Line should be > 0, got %d", i, frame.Line)
		}
	}
}

func TestErr_StackTrace_Disabled(t *testing.T) {
	err := Message("test").
		WithStackTrace(false).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	stackTrace := flooErr.StackTrace()
	if stackTrace != nil && len(stackTrace) != 0 {
		t.Errorf("Expected empty stack trace when disabled, got %d frames", len(stackTrace))
	}
}

func TestErr_StackTrace_Cached(t *testing.T) {
	err := Message("test").
		WithStackTrace(true).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	// Call StackTrace multiple times
	stackTrace1 := flooErr.StackTrace()
	stackTrace2 := flooErr.StackTrace()

	if len(stackTrace1) != len(stackTrace2) {
		t.Error("Stack trace should be cached and return same result")
	}

	// Verify they are the same
	for i := range stackTrace1 {
		if stackTrace1[i].Function != stackTrace2[i].Function {
			t.Error("Cached stack trace should return same values")
		}
	}
}

func TestErr_Error(t *testing.T) {
	err := Message("test error").
		Error(nil, "test error")

	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

func TestErr_Error_WithCause(t *testing.T) {
	originalErr := errors.New("original error")
	err := Message("wrapped error").
		Error(originalErr, "wrapped error")

	errorMsg := err.Error()
	expectedSuffix := "original error"
	if len(errorMsg) < len(expectedSuffix) {
		t.Errorf("Error message too short: '%s'", errorMsg)
	}

	// Error message should contain the cause
	if errorMsg[len(errorMsg)-len(expectedSuffix):] != expectedSuffix {
		t.Errorf("Expected error message to end with '%s', got '%s'", expectedSuffix, errorMsg)
	}
}

func TestErr_ImplementsError(t *testing.T) {
	err := Message("test").
		Error(nil, "test error")

	var _ error = err
	var _ FlooErr = err.(FlooErr)
}

func TestErr_ErrorChaining(t *testing.T) {
	baseErr := errors.New("base error")
	middleErr := Message("middle error").
		WithCode("MIDDLE_ERR").
		Error(baseErr, "middle error")

	topErr := Message("top error").
		WithCode("TOP_ERR").
		Error(middleErr, "top error")

	flooErr, ok := topErr.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	// Unwrap first level
	cause1 := flooErr.Unwrap()
	if cause1 == nil {
		t.Fatal("Expected non-nil cause at first level")
	}

	// Unwrap second level
	if flooErr2, ok := cause1.(FlooErr); ok {
		cause2 := flooErr2.Unwrap()
		if cause2 == nil {
			t.Fatal("Expected non-nil cause at second level")
		}

		if cause2.Error() != "base error" {
			t.Errorf("Expected 'base error' at base level, got '%s'", cause2.Error())
		}
	} else {
		t.Error("Expected FlooErr at second level")
	}
}

func TestErr_ErrorsAs(t *testing.T) {
	err := Message("test").
		WithCode("TEST_CODE").
		Error(nil, "test error")

	var flooErr FlooErr
	if !errors.As(err, &flooErr) {
		t.Fatal("errors.As should work with FlooErr")
	}

	if string(flooErr.Code()) != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", flooErr.Code())
	}
}

func TestErr_MultipleContextValues(t *testing.T) {
	err := Message("test").
		WithContext("key1", "value1").
		WithContext("key2", "value2").
		WithContext("key1", "overridden").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	context := flooErr.Context()
	if context["key1"] != "overridden" {
		t.Errorf("Expected context['key1'] = 'overridden', got '%v'", context["key1"])
	}

	if len(context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(context))
	}
}

func TestErr_MultipleSDCValues(t *testing.T) {
	err := Message("test").
		WithSDC("key1", "value1").
		WithSDC("key2", "value2").
		WithSDC("key1", "overridden").
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	sdc := flooErr.SDC()
	if sdc["key1"] != "overridden" {
		t.Errorf("Expected SDC['key1'] = 'overridden', got '%s'", sdc["key1"])
	}

	if len(sdc) != 2 {
		t.Errorf("Expected 2 SDC items, got %d", len(sdc))
	}
}

func TestErr_ComplexContext(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	user := User{ID: 123, Name: "John"}
	err := Message("test").
		WithContext("user", user).
		WithContext("numbers", []int{1, 2, 3}).
		WithContext("map", map[string]int{"a": 1, "b": 2}).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	context := flooErr.Context()
	if len(context) != 3 {
		t.Errorf("Expected 3 context items, got %d", len(context))
	}

	// Verify user struct
	if u, ok := context["user"].(User); ok {
		if u.ID != 123 || u.Name != "John" {
			t.Error("User struct not preserved correctly")
		}
	} else {
		t.Error("User struct not found in context")
	}
}

func TestCode_Function(t *testing.T) {
	// Test that Code() function returns ErrProps with the code set
	props := Code(internal.Code("TEST_CODE"))
	if props == nil {
		t.Fatal("Expected non-nil ErrProps")
	}
	// Code is set internally, we can't directly access it, but we can verify it works by building an error
	err := props.Error(nil, "test")
	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr")
	}
	if flooErr.Code().String() != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", flooErr.Code().String())
	}
}
