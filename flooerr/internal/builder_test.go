package internal

import (
	"errors"
	"testing"
)

func TestCreate(t *testing.T) {
	props := Create()
	if props == nil {
		t.Fatal("Create() returned nil")
	}

	if props.context == nil {
		t.Error("Expected non-nil context map")
	}

	if props.sdc == nil {
		t.Error("Expected non-nil SDC map")
	}

	if !props.withStackTrace {
		t.Error("Expected withStackTrace to be true by default")
	}
}

func TestErrProps_WithMessage(t *testing.T) {
	props := Create().WithMessage("test message")
	if props.message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", props.message)
	}
}

func TestErrProps_WithMessage_Chaining(t *testing.T) {
	props := Create().
		WithMessage("first message").
		WithMessage("second message")

	if props.message != "second message" {
		t.Errorf("Expected message 'second message', got '%s'", props.message)
	}
}

func TestErrProps_WithCode(t *testing.T) {
	props := Create().WithCode("TEST_CODE")
	if props.code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", props.code)
	}
}

func TestErrProps_WithCode_Chaining(t *testing.T) {
	props := Create().
		WithCode("FIRST_CODE").
		WithCode("SECOND_CODE")

	if props.code != "SECOND_CODE" {
		t.Errorf("Expected code 'SECOND_CODE', got '%s'", props.code)
	}
}

func TestErrProps_WithStackTrace(t *testing.T) {
	props := Create().WithStackTrace(false)
	if props.withStackTrace {
		t.Error("Expected withStackTrace to be false")
	}

	props = props.WithStackTrace(true)
	if !props.withStackTrace {
		t.Error("Expected withStackTrace to be true")
	}
}

func TestErrProps_WithContext(t *testing.T) {
	props := Create().
		WithContext("key1", "value1").
		WithContext("key2", 123).
		WithContext("key3", true)

	if len(props.context) != 3 {
		t.Errorf("Expected 3 context items, got %d", len(props.context))
	}

	if props.context["key1"] != "value1" {
		t.Errorf("Expected context['key1'] = 'value1', got '%v'", props.context["key1"])
	}

	if props.context["key2"] != 123 {
		t.Errorf("Expected context['key2'] = 123, got '%v'", props.context["key2"])
	}

	if props.context["key3"] != true {
		t.Errorf("Expected context['key3'] = true, got '%v'", props.context["key3"])
	}
}

func TestErrProps_WithContext_Override(t *testing.T) {
	props := Create().
		WithContext("key", "value1").
		WithContext("key", "value2")

	if props.context["key"] != "value2" {
		t.Errorf("Expected context['key'] = 'value2', got '%v'", props.context["key"])
	}

	if len(props.context) != 1 {
		t.Errorf("Expected 1 context item, got %d", len(props.context))
	}
}

func TestErrProps_WithSDC(t *testing.T) {
	props := Create().
		WithSDC("trace_id", "trace_123").
		WithSDC("span_id", "span_456")

	if len(props.sdc) != 2 {
		t.Errorf("Expected 2 SDC items, got %d", len(props.sdc))
	}

	if props.sdc["trace_id"] != "trace_123" {
		t.Errorf("Expected SDC['trace_id'] = 'trace_123', got '%s'", props.sdc["trace_id"])
	}

	if props.sdc["span_id"] != "span_456" {
		t.Errorf("Expected SDC['span_id'] = 'span_456', got '%s'", props.sdc["span_id"])
	}
}

func TestErrProps_WithSDC_Override(t *testing.T) {
	props := Create().
		WithSDC("key", "value1").
		WithSDC("key", "value2")

	if props.sdc["key"] != "value2" {
		t.Errorf("Expected SDC['key'] = 'value2', got '%s'", props.sdc["key"])
	}

	if len(props.sdc) != 1 {
		t.Errorf("Expected 1 SDC item, got %d", len(props.sdc))
	}
}

func TestErrProps_Build_WithoutCause(t *testing.T) {
	err := Create().
		WithMessage("test error").
		WithCode("TEST_CODE").
		Error(nil, "default message")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

func TestErrProps_Build_WithCause(t *testing.T) {
	originalErr := errors.New("original error")
	err := Create().
		WithMessage("wrapped error").
		WithCode("WRAP_CODE").
		Error(originalErr, "wrapped error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	// Error should contain the cause
	errorMsg := err.Error()
	if len(errorMsg) == 0 {
		t.Error("Expected non-empty error message")
	}
}

func TestErrProps_Build_DefaultMessage(t *testing.T) {
	err := Create().
		WithCode("TEST_CODE").
		Error(nil, "default message")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	// When message is empty, Error uses the fallback message
	if err.Error() != "default message" {
		t.Errorf("Expected 'default message', got '%s'", err.Error())
	}
}

func TestErrProps_Build_WithContext(t *testing.T) {
	err := Create().
		WithMessage("test").
		WithContext("key1", "value1").
		WithContext("key2", 123).
		Error(nil, "test error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}
}

func TestErrProps_Build_WithSDC(t *testing.T) {
	err := Create().
		WithMessage("test").
		WithSDC("trace_id", "trace_123").
		WithSDC("span_id", "span_456").
		Error(nil, "test error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}
}

func TestErrProps_Build_WithStackTrace(t *testing.T) {
	err := Create().
		WithMessage("test").
		WithStackTrace(true).
		Error(nil, "test error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}
}

func TestErrProps_Build_WithoutStackTrace(t *testing.T) {
	err := Create().
		WithMessage("test").
		WithStackTrace(false).
		Error(nil, "test error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}
}

func TestErrProps_Build_AllProperties(t *testing.T) {
	originalErr := errors.New("original error")
	err := Create().
		WithMessage("test error").
		WithCode("TEST_CODE").
		WithContext("key1", "value1").
		WithContext("key2", 123).
		WithSDC("trace_id", "trace_123").
		WithSDC("span_id", "span_456").
		WithStackTrace(true).
		Error(originalErr, "default message")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	errorMsg := err.Error()
	if len(errorMsg) == 0 {
		t.Error("Expected non-empty error message")
	}
}

func TestErrProps_Build_Chaining(t *testing.T) {
	err := Create().
		WithMessage("first").
		WithMessage("second").
		WithCode("CODE1").
		WithCode("CODE2").
		WithContext("key1", "value1").
		WithContext("key2", "value2").
		WithSDC("sdc1", "value1").
		WithSDC("sdc2", "value2").
		Error(nil, "default")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	// Verify last values are used
	if err.Error() != "second" {
		t.Errorf("Expected 'second', got '%s'", err.Error())
	}
}

func TestSetBuildErrFunc(t *testing.T) {
	originalFunc := buildErrFunc

	// Set a custom builder function
	customBuilder := func(
		message string,
		errMessage string,
		code Code,
		cause error,
		stackTracePTR []uintptr,
		context map[string]any,
		sdc map[string]string,
	) error {
		return errors.New("custom error")
	}

	SetBuildErrFunc(customBuilder)

	if buildErrFunc == nil {
		t.Fatal("Expected non-nil buildErrFunc")
	}

	// Test that it's used
	err := Create().
		WithMessage("test").
		Error(nil, "test error")

	if err.Error() != "custom error" {
		t.Errorf("Expected 'custom error', got '%s'", err.Error())
	}

	// Restore original function
	SetBuildErrFunc(originalFunc)
}

func TestBuildErrFunc_Nil(t *testing.T) {
	originalFunc := buildErrFunc
	SetBuildErrFunc(nil)

	// Should fallback to simpleError
	err := Create().
		WithMessage("test").
		Error(nil, "test error")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	// Restore original function
	SetBuildErrFunc(originalFunc)
}

func TestBuildErrFunc_WithCause(t *testing.T) {
	originalFunc := buildErrFunc
	SetBuildErrFunc(nil)

	originalErr := errors.New("original")
	err := Create().
		WithMessage("test").
		Error(originalErr, "test error")

	if err != originalErr {
		t.Error("Expected original error to be returned when builder is nil and cause exists")
	}

	// Restore original function
	SetBuildErrFunc(originalFunc)
}

func TestCode_Type(t *testing.T) {
	code := Code("TEST_CODE")
	if string(code) != "TEST_CODE" {
		t.Errorf("Expected 'TEST_CODE', got '%s'", code)
	}
}
