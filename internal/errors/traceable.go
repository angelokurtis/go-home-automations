package errors

// Traceable is an interface that errors can implement to provide additional
// information about their cause and stack trace.
type Traceable interface {
	// Unwrap returns the underlying cause of the error that this error wraps.
	Unwrap() error

	// Stack returns the call stack for this error.
	Stack() *Stack

	// String returns a string representation of the call stack.
	String() string
}
