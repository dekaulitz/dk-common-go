package flooerr

import (
	"strings"
	"testing"
)

func TestStacktrace_String(t *testing.T) {
	st := stacktrace{
		Function: "main.testFunction",
		File:     "/path/to/file.go",
		Line:     42,
	}

	expected := "main.testFunction:/path/to/file.go:42"
	actual := st.String()

	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func TestStacktrace_String_Empty(t *testing.T) {
	st := stacktrace{
		Function: "",
		File:     "",
		Line:     0,
	}

	result := st.String()
	if !strings.Contains(result, ":") {
		t.Error("String() should contain colons even with empty values")
	}
}

func TestStacktrace_String_RealStack(t *testing.T) {
	// Create an error with stack trace enabled
	err := Message("test").
		WithStackTrace(true).
		Error(nil, "test error")

	flooErr, ok := err.(FlooErr)
	if !ok {
		t.Fatal("Expected FlooErr interface")
	}

	stackTrace := flooErr.StackTrace()
	if len(stackTrace) == 0 {
		t.Skip("Stack trace is empty, skipping test")
	}

	// Test String() method on real stack trace
	for i, frame := range stackTrace {
		str := frame.String()
		if str == "" {
			t.Errorf("Stack frame %d: String() should not be empty", i)
		}

		// Should contain colons (format: Function:File:Line)
		if !strings.Contains(str, ":") {
			t.Errorf("Stack frame %d: String() should contain colons, got '%s'", i, str)
		}
	}
}

func TestStacktrace_Fields(t *testing.T) {
	st := stacktrace{
		Function: "package.Function",
		File:     "/absolute/path/file.go",
		Line:     100,
	}

	if st.Function != "package.Function" {
		t.Errorf("Expected Function 'package.Function', got '%s'", st.Function)
	}

	if st.File != "/absolute/path/file.go" {
		t.Errorf("Expected File '/absolute/path/file.go', got '%s'", st.File)
	}

	if st.Line != 100 {
		t.Errorf("Expected Line 100, got %d", st.Line)
	}
}
