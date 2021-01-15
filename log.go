package main

import (
	"log"
	"os"
)

// NewLogger creates the logger for the app.
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
