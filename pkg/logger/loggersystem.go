package logger

// Log is the abstraction for sending
// and collecting log data from the system
type Log interface {
	Info(message string) string
	Error(message string) string
}

// TODO: Think how the users interact with the application as a whole.
// i.e. rethink use cases. Maybe have a more general "Dev User"

// UserLog is the log data sent back to
// users to know their request status
type UserLog struct {
	Log Log
}

// DevLog is log data useful for debugging
// and development of the system. Used by developers.
type DevLog struct {
	UserLog
}
