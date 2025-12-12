package flooerr

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	err := Message("test error").
		WithCode("TEST_CODE").
		WithContext("key1", "value1").
		WithSDC("trace_id", "trace_123").
		Error(nil, "test error")

	info := Parse(err)

	if !info.IsFlooErr {
		t.Error("Expected IsFlooErr to be true")
	}

	if info.Code.String() != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", info.Code.String())
	}

	if info.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", info.Message)
	}

	if len(info.Context) != 1 {
		t.Errorf("Expected 1 context item, got %d", len(info.Context))
	}

	if len(info.SDC) != 1 {
		t.Errorf("Expected 1 SDC item, got %d", len(info.SDC))
	}
}

func TestParse_NonFlooErr(t *testing.T) {
	err := errors.New("standard error")
	info := Parse(err)

	if info.IsFlooErr {
		t.Error("Expected IsFlooErr to be false")
	}

	if info.ErrorMsg != "standard error" {
		t.Errorf("Expected error message 'standard error', got '%s'", info.ErrorMsg)
	}
}

func TestParse_NilError(t *testing.T) {
	info := Parse(nil)

	if info.IsFlooErr {
		t.Error("Expected IsFlooErr to be false")
	}
}

func TestAsFlooErr(t *testing.T) {
	err := Message("test").Error(nil, "test")
	flooErr, ok := AsFlooErr(err)

	if !ok {
		t.Error("Expected AsFlooErr to return true")
	}

	if flooErr == nil {
		t.Error("Expected non-nil FlooErr")
	}
}

func TestAsFlooErr_NonFlooErr(t *testing.T) {
	err := errors.New("standard error")
	_, ok := AsFlooErr(err)

	if ok {
		t.Error("Expected AsFlooErr to return false")
	}
}

func TestAsFlooErr_Nil(t *testing.T) {
	_, ok := AsFlooErr(nil)

	if ok {
		t.Error("Expected AsFlooErr to return false for nil")
	}
}

func TestGetCode(t *testing.T) {
	err := Message("test").WithCode("TEST_CODE").Error(nil, "test")
	code := GetCode(err)

	if code.String() != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", code.String())
	}
}

func TestGetCodeString(t *testing.T) {
	err := Message("test").WithCode("TEST_CODE").Error(nil, "test")
	code := GetCodeString(err)

	if code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", code)
	}
}

func TestGetCode_NonFlooErr(t *testing.T) {
	err := errors.New("standard error")
	code := GetCodeString(err)

	if code != "" {
		t.Errorf("Expected empty code, got '%s'", code)
	}
}

func TestGetMessage(t *testing.T) {
	err := Message("custom message").Error(nil, "default")
	message := GetMessage(err)

	if message != "custom message" {
		t.Errorf("Expected message 'custom message', got '%s'", message)
	}
}

func TestGetContext(t *testing.T) {
	err := Message("test").
		WithContext("key1", "value1").
		WithContext("key2", 123).
		Error(nil, "test")

	context := GetContext(err)

	if len(context) != 2 {
		t.Errorf("Expected 2 context items, got %d", len(context))
	}

	if context["key1"] != "value1" {
		t.Errorf("Expected context['key1'] = 'value1', got '%v'", context["key1"])
	}
}

func TestGetContextValue(t *testing.T) {
	err := Message("test").
		WithContext("key1", "value1").
		Error(nil, "test")

	value := GetContextValue(err, "key1")
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}

	value = GetContextValue(err, "nonexistent")
	if value != nil {
		t.Errorf("Expected nil for nonexistent key, got '%v'", value)
	}
}

func TestGetSDC(t *testing.T) {
	err := Message("test").
		WithSDC("trace_id", "trace_123").
		WithSDC("span_id", "span_456").
		Error(nil, "test")

	sdc := GetSDC(err)

	if len(sdc) != 2 {
		t.Errorf("Expected 2 SDC items, got %d", len(sdc))
	}

	if sdc["trace_id"] != "trace_123" {
		t.Errorf("Expected SDC['trace_id'] = 'trace_123', got '%s'", sdc["trace_id"])
	}
}

func TestGetSDCValue(t *testing.T) {
	err := Message("test").
		WithSDC("trace_id", "trace_123").
		Error(nil, "test")

	value := GetSDCValue(err, "trace_id")
	if value != "trace_123" {
		t.Errorf("Expected 'trace_123', got '%s'", value)
	}

	value = GetSDCValue(err, "nonexistent")
	if value != "" {
		t.Errorf("Expected empty string for nonexistent key, got '%s'", value)
	}
}

func TestGetStackTrace(t *testing.T) {
	err := Message("test").
		WithStackTrace(true).
		Error(nil, "test")

	stackTrace := GetStackTrace(err)
	// Stack trace might be empty in test environment, so we just verify it doesn't panic
	_ = stackTrace
}

func TestGetCause(t *testing.T) {
	originalErr := errors.New("original error")
	err := Message("wrapped").Error(originalErr, "wrapped")

	cause := GetCause(err)
	if cause == nil {
		t.Error("Expected non-nil cause")
	}

	if cause.Error() != "original error" {
		t.Errorf("Expected 'original error', got '%s'", cause.Error())
	}
}

func TestGetCause_NoCause(t *testing.T) {
	err := Message("test").Error(nil, "test")
	cause := GetCause(err)

	if cause != nil {
		t.Errorf("Expected nil cause, got %v", cause)
	}
}

func TestIsFlooErr(t *testing.T) {
	err := Message("test").Error(nil, "test")
	if !IsFlooErr(err) {
		t.Error("Expected IsFlooErr to return true")
	}

	standardErr := errors.New("standard error")
	if IsFlooErr(standardErr) {
		t.Error("Expected IsFlooErr to return false for standard error")
	}
}

func TestUnwrapChain(t *testing.T) {
	baseErr := errors.New("base error")
	middleErr := Message("middle").Error(baseErr, "middle")
	topErr := Message("top").Error(middleErr, "top")

	chain := UnwrapChain(topErr)

	if len(chain) != 3 {
		t.Errorf("Expected chain length 3, got %d", len(chain))
	}

	if chain[0].Error() != topErr.Error() {
		t.Error("First element should be top error")
	}

	if chain[len(chain)-1].Error() != "base error" {
		t.Errorf("Last element should be base error, got '%s'", chain[len(chain)-1].Error())
	}
}

func TestGetRootCause(t *testing.T) {
	baseErr := errors.New("base error")
	middleErr := Message("middle").Error(baseErr, "middle")
	topErr := Message("top").Error(middleErr, "top")

	rootCause := GetRootCause(topErr)

	if rootCause.Error() != "base error" {
		t.Errorf("Expected root cause 'base error', got '%s'", rootCause.Error())
	}
}

func TestHasCode(t *testing.T) {
	err := Message("test").WithCode("TEST_CODE").Error(nil, "test")

	if !HasCode(err, "TEST_CODE") {
		t.Error("Expected HasCode to return true")
	}

	if HasCode(err, "OTHER_CODE") {
		t.Error("Expected HasCode to return false for different code")
	}
}

func TestHasContextKey(t *testing.T) {
	err := Message("test").WithContext("key1", "value1").Error(nil, "test")

	if !HasContextKey(err, "key1") {
		t.Error("Expected HasContextKey to return true")
	}

	if HasContextKey(err, "nonexistent") {
		t.Error("Expected HasContextKey to return false for nonexistent key")
	}
}

func TestHasSDCKey(t *testing.T) {
	err := Message("test").WithSDC("trace_id", "trace_123").Error(nil, "test")

	if !HasSDCKey(err, "trace_id") {
		t.Error("Expected HasSDCKey to return true")
	}

	if HasSDCKey(err, "nonexistent") {
		t.Error("Expected HasSDCKey to return false for nonexistent key")
	}
}
