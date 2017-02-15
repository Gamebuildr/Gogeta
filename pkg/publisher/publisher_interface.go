package publisher

import "github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"

// Application is the interface to specify a notification service
type Application interface {
	PublishMessage(msg Message) (string, error)
}

// Service is the base system for creating unique messaging services
type Service struct {
	Application Application
	Log         logger.Log
}
