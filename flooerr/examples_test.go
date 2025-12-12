package flooerr_test

import (
	"core-common-go/flooerr"
	"errors"
	"fmt"
)

// ExampleMessage demonstrates creating a simple error with a message
func ExampleMessage() {
	err := flooerr.Message("User not found").
		WithCode("USER_NOT_FOUND").
		Error(nil, "Failed to retrieve user")

	fmt.Println(err.Error())
	// Output: User not found
}

// ExampleMessage_context demonstrates adding context to an error
func ExampleMessage_context() {
	err := flooerr.Message("Database connection failed").
		WithCode("DB_CONN_ERR").
		WithContext("host", "localhost").
		WithContext("port", 5432).
		WithContext("database", "mydb").
		Error(nil, "Unable to connect to database")

	if flooErr, ok := err.(flooerr.FlooErr); ok {
		fmt.Printf("Code: %s\n", flooErr.Code())
		fmt.Printf("Message: %s\n", flooErr.Message())
		fmt.Printf("Context: %v\n", flooErr.Context())
	}
}

// ExampleMessage_sdc demonstrates using Structured Diagnostic Context
func ExampleMessage_sdc() {
	err := flooerr.Message("Payment processing failed").
		WithCode("PAYMENT_ERR").
		WithSDC("transaction_id", "txn_12345").
		WithSDC("user_id", "user_67890").
		Error(nil, "Payment could not be processed")

	if flooErr, ok := err.(flooerr.FlooErr); ok {
		fmt.Printf("SDC: %v\n", flooErr.SDC())
	}
}

// ExampleMessage_wrapping demonstrates wrapping an existing error
func ExampleMessage_wrapping() {
	originalErr := errors.New("file not found: config.json")
	err := flooerr.Message("Configuration error").
		WithCode("CONFIG_ERR").
		Error(originalErr, "Failed to load configuration")

	if flooErr, ok := err.(flooerr.FlooErr); ok {
		fmt.Printf("Error: %s\n", flooErr.Error())
		if cause := flooErr.Unwrap(); cause != nil {
			fmt.Printf("Caused by: %v\n", cause)
		}
	}
}

// ExampleFlooErr_stackTrace demonstrates accessing stack trace information
func ExampleFlooErr_stackTrace() {
	err := flooerr.Message("Operation failed").
		WithCode("OP_ERR").
		Error(nil, "Failed to perform operation")

	if flooErr, ok := err.(flooerr.FlooErr); ok {
		stackTrace := flooErr.StackTrace()
		if len(stackTrace) > 0 {
			fmt.Printf("Stack trace frames: %d\n", len(stackTrace))
			// Print first frame as example
			if len(stackTrace) > 0 {
				frame := stackTrace[0]
				fmt.Printf("First frame: %s\n", frame.String())
			}
		}
	}
}

// ExampleFlooErr_errorChain demonstrates error chaining and unwrapping
func ExampleFlooErr_errorChain() {
	// Create a chain of errors
	baseErr := errors.New("database connection timeout")

	middleErr := flooerr.Message("Query execution failed").
		WithCode("QUERY_ERR").
		WithContext("query", "SELECT * FROM users").
		Error(baseErr, "Failed to execute query")

	topErr := flooerr.Message("User retrieval failed").
		WithCode("USER_RETRIEVAL_ERR").
		WithContext("user_id", "12345").
		Error(middleErr, "Could not retrieve user")

	// Unwrap the error chain
	currentErr := topErr
	depth := 0
	for currentErr != nil {
		if flooErr, ok := currentErr.(flooerr.FlooErr); ok {
			fmt.Printf("Level %d: %s (Code: %s)\n", depth, flooErr.Message(), flooErr.Code())
			currentErr = flooErr.Unwrap()
		} else {
			fmt.Printf("Level %d: %v\n", depth, currentErr)
			break
		}
		depth++
	}
}

// ExampleMessage_noStackTrace demonstrates disabling stack trace for performance
func ExampleMessage_noStackTrace() {
	err := flooerr.Message("High-frequency error").
		WithCode("HIGH_FREQ_ERR").
		WithStackTrace(false). // Disable stack trace capture
		Error(nil, "Error occurred")

	if flooErr, ok := err.(flooerr.FlooErr); ok {
		stackTrace := flooErr.StackTrace()
		fmt.Printf("Stack trace enabled: %v\n", stackTrace != nil && len(stackTrace) > 0)
	}
}

// ExampleFlooErr_errorChecking demonstrates checking error types and codes
func ExampleFlooErr_errorChecking() {
	err := flooerr.Message("User not found").
		WithCode("USER_NOT_FOUND").
		Error(nil, "User does not exist")

	var flooErr flooerr.FlooErr
	if errors.As(err, &flooErr) {
		code := flooErr.Code()
		if string(code) == "USER_NOT_FOUND" {
			fmt.Println("Handling user not found error")
		}
	}
}
