package schema

import "log"

type Logger struct {
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
}
