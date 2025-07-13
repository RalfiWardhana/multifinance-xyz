package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func InitLogger(level string) {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string, args ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(msg, args...)
	}
}
