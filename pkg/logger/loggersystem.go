package logger

// Log is the abstraction for sending
// and collecting log data from the system
type Log interface {
	Info(message string) string
	Error(message string) string
}
