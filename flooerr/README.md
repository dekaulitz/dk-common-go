# FlooErr - Error Wrapper Library for Go

FlooErr is a comprehensive error handling library for Go that provides structured error wrapping with stack traces, error codes, context, and SDC (Structured Diagnostic Context) support. It's designed to help developers create more informative and debuggable error messages in their Go applications.

## Features

- **Structured Error Wrapping**: Wrap errors with additional metadata including codes, messages, and context
- **Stack Trace Support**: Automatic stack trace capture for better debugging
- **Error Chaining**: Support for error cause chaining using `Unwrap()`
- **Context Information**: Attach arbitrary context data to errors
- **SDC Support**: Structured Diagnostic Context for logging and monitoring
- **Type-Safe Error Codes**: Strongly typed error codes
- **Fluent API**: Builder pattern for easy error construction

## Installation

```bash
go get core-common-go/flooerr
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "core-common-go/flooerr"
)

func main() {
    // Create a simple error
    err := flooerr.Message("Something went wrong").
        WithCode("ERR_001").
        Build(nil, "Failed to process request")
    
    if err != nil {
        fmt.Println(err.Error())
    }
}
```

### Creating Errors with Builder Pattern

```go
import "core-common-go/flooerr"

// Simple error with message
err := flooerr.Message("User not found").
    WithCode("USER_NOT_FOUND").
    Build(nil, "Failed to retrieve user")

// Error with context
err := flooerr.Message("Database connection failed").
    WithCode("DB_CONN_ERR").
    WithContext("host", "localhost").
    WithContext("port", 5432).
    WithContext("database", "mydb").
    Build(nil, "Unable to connect to database")

// Error with SDC (Structured Diagnostic Context)
err := flooerr.Message("Payment processing failed").
    WithCode("PAYMENT_ERR").
    WithSDC("transaction_id", "txn_12345").
    WithSDC("user_id", "user_67890").
    Build(nil, "Payment could not be processed")

// Error wrapping an existing error
originalErr := fmt.Errorf("file not found: config.json")
err := flooerr.Message("Configuration error").
    WithCode("CONFIG_ERR").
    Build(originalErr, "Failed to load configuration")
```

### Working with FlooErr Interface

```go
import "core-common-go/flooerr"

func handleError(err error) {
    if flooErr, ok := err.(flooerr.FlooErr); ok {
        // Access error code
        code := flooErr.Code()
        fmt.Printf("Error Code: %s\n", code)
        
        // Access error message
        message := flooErr.Message()
        fmt.Printf("Message: %s\n", message)
        
        // Access stack trace
        stackTrace := flooErr.StackTrace()
        for _, frame := range stackTrace {
            fmt.Printf("  %s:%s:%d\n", frame.Function, frame.File, frame.Line)
        }
        
        // Access context
        context := flooErr.Context()
        for key, value := range context {
            fmt.Printf("  %s: %v\n", key, value)
        }
        
        // Access SDC
        sdc := flooErr.SDC()
        for key, value := range sdc {
            fmt.Printf("  %s: %s\n", key, value)
        }
        
        // Unwrap the underlying error
        if cause := flooErr.Unwrap(); cause != nil {
            fmt.Printf("Caused by: %v\n", cause)
        }
    }
}
```

## API Reference

### FlooErr Interface

The `FlooErr` interface extends Go's standard `error` interface with additional methods:

```go
type FlooErr interface {
    error
    Code() internal.Code      // Returns the error code
    Message() string          // Returns the error message
    StackTrace() []stacktrace // Returns the stack trace frames
    Unwrap() error            // Returns the underlying cause error
    Context() map[string]any  // Returns the context map
    SDC() map[string]string  // Returns the SDC map
}
```

### Builder Methods

#### `Message(msg string) *ErrProps`

Creates a new error builder with the specified message.

```go
builder := flooerr.Message("Error message")
```

#### `WithCode(code string) *ErrProps`

Sets the error code for the error.

```go
builder := flooerr.Message("Error").
    WithCode("ERR_001")
```

#### `WithMessage(message string) *ErrProps`

Sets or updates the error message.

```go
builder := flooerr.Message("Initial message").
    WithMessage("Updated message")
```

#### `WithStackTrace(enable bool) *ErrProps`

Enables or disables stack trace capture. Default is `true`.

```go
builder := flooerr.Message("Error").
    WithStackTrace(false) // Disable stack trace
```

#### `WithContext(key string, value any) *ErrProps`

Adds a key-value pair to the error context. Can be called multiple times to add multiple context values.

```go
builder := flooerr.Message("Error").
    WithContext("user_id", 12345).
    WithContext("request_id", "req_abc")
```

#### `WithSDC(key string, value string) *ErrProps`

Adds a key-value pair to the Structured Diagnostic Context (SDC). SDC values must be strings and are typically used for logging and monitoring.

```go
builder := flooerr.Message("Error").
    WithSDC("trace_id", "trace_123").
    WithSDC("span_id", "span_456")
```

#### `Build(cause error, message string) error`

Builds and returns the final error. If `cause` is provided, it will be wrapped. The `message` parameter is used as a fallback if no message was set via `WithMessage()`.

```go
err := builder.Build(nil, "Default error message")
// or
err := builder.Build(originalErr, "Wrapping error message")
```

### Stack Trace

The `StackTrace()` method returns a slice of `stacktrace` structs, each containing:

```go
type stacktrace struct {
    Function string // Function name
    File     string // File path
    Line     int    // Line number
}
```

## Advanced Usage

### Error Code Types

Error codes are strongly typed. You can define constants for your error codes:

```go
const (
    ErrCodeUserNotFound    = flooerr.Code("USER_NOT_FOUND")
    ErrCodeInvalidInput    = flooerr.Code("INVALID_INPUT")
    ErrCodeDatabaseError   = flooerr.Code("DB_ERROR")
    ErrCodeNetworkError    = flooerr.Code("NETWORK_ERROR")
)

// Usage
err := flooerr.Message("User not found").
    WithCode(string(ErrCodeUserNotFound)).
    Build(nil, "Failed to find user")
```

### Error Propagation

FlooErr supports error wrapping and unwrapping, making it compatible with Go's error handling patterns:

```go
func processUser(userID string) error {
    user, err := getUser(userID)
    if err != nil {
        return flooerr.Message("Failed to process user").
            WithCode("PROCESS_USER_ERR").
            WithContext("user_id", userID).
            Build(err, "Error processing user")
    }
    // ... process user
    return nil
}

func getUser(userID string) (*User, error) {
    // ... database query
    if notFound {
        return nil, flooerr.Message("User not found").
            WithCode("USER_NOT_FOUND").
            Build(nil, "User does not exist")
    }
    return user, nil
}
```

### Checking Error Types

```go
import (
    "errors"
    "core-common-go/flooerr"
)

func handleError(err error) {
    var flooErr flooerr.FlooErr
    if errors.As(err, &flooErr) {
        // Handle FlooErr
        if flooErr.Code() == "USER_NOT_FOUND" {
            // Handle specific error code
        }
    }
}
```

### Disabling Stack Traces for Performance

If stack trace capture is not needed (e.g., in production for performance reasons), you can disable it:

```go
err := flooerr.Message("Error").
    WithStackTrace(false).
    Build(nil, "Error message")
```

## Best Practices

1. **Use Meaningful Error Codes**: Define constants for error codes and use them consistently across your application.

2. **Add Context**: Include relevant context information that will help with debugging:
   ```go
   err := flooerr.Message("Operation failed").
       WithContext("operation", "create_user").
       WithContext("input", userInput).
       Build(originalErr, "Failed to create user")
   ```

3. **Use SDC for Logging**: Use SDC for values that should be included in structured logs:
   ```go
   err := flooerr.Message("Request failed").
       WithSDC("request_id", requestID).
       WithSDC("user_id", userID).
       Build(nil, "Request processing failed")
   ```

4. **Wrap Existing Errors**: Always wrap underlying errors to preserve the error chain:
   ```go
   if err != nil {
       return flooerr.Message("High-level error").
           WithCode("HIGH_LEVEL_ERR").
           Build(err, "Operation failed")
   }
   ```

5. **Disable Stack Traces When Not Needed**: In high-performance scenarios, consider disabling stack traces:
   ```go
   err := flooerr.Message("Error").
       WithStackTrace(false).
       Build(nil, "Error message")
   ```

## Examples

### Example: HTTP Handler Error Handling

```go
package main

import (
    "encoding/json"
    "net/http"
    "core-common-go/flooerr"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        err := flooerr.Message("Invalid request").
            WithCode("INVALID_REQUEST").
            WithContext("missing_param", "id").
            Build(nil, "User ID is required")
        
        writeError(w, err, http.StatusBadRequest)
        return
    }
    
    user, err := getUser(userID)
    if err != nil {
        writeError(w, err, http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
    w.WriteHeader(statusCode)
    
    if flooErr, ok := err.(flooerr.FlooErr); ok {
        response := map[string]interface{}{
            "error": map[string]interface{}{
                "code":    flooErr.Code(),
                "message": flooErr.Message(),
                "context": flooErr.Context(),
            },
        }
        json.NewEncoder(w).Encode(response)
    } else {
        json.NewEncoder(w).Encode(map[string]string{
            "error": err.Error(),
        })
    }
}
```

### Example: Database Error Handling

```go
package main

import (
    "database/sql"
    "core-common-go/flooerr"
)

func getUserFromDB(db *sql.DB, userID string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = $1"
    row := db.QueryRow(query, userID)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, flooerr.Message("User not found").
                WithCode("USER_NOT_FOUND").
                WithContext("user_id", userID).
                Build(err, "User does not exist in database")
        }
        
        return nil, flooerr.Message("Database query failed").
            WithCode("DB_QUERY_ERR").
            WithContext("query", query).
            WithContext("user_id", userID).
            Build(err, "Failed to query user from database")
    }
    
    return &user, nil
}
```

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
