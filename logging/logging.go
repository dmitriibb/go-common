package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type logLevel string

const (
	DEBUG logLevel = "DEBUG"
	INFO  logLevel = "INFO"
	WARN  logLevel = "WARN"
	ERROR logLevel = "ERROR"
)

var loggerDebug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
var loggerInfo = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var loggerWarn = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
var loggerError = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

func LogInfo(message string) {
	loggerInfo.Println(message)
}
func logMessage(loggerName string, prefix logLevel, messageArgs ...any) {
	message := argsToMessage(messageArgs...)
	switch prefix {
	case DEBUG:
		loggerDebug.Println(fmt.Sprintf("%v - %v", loggerName, message))
	case INFO:
		loggerInfo.Println(fmt.Sprintf("%v - %v", loggerName, message))
	case WARN:
		loggerWarn.Println(fmt.Sprintf("%v - %v", loggerName, message))
	case ERROR:
		loggerError.Println(fmt.Sprintf("%v - %v", loggerName, message))
	}
}

func argsToMessage(messageArgs ...any) string {
	if len(messageArgs) > 1 {
		return fmt.Sprintf(messageArgs[0].(string), messageArgs[1:]...)
	} else {
		return messageArgs[0].(string)
	}
}

type customLogger struct {
	loggerName string
}

type Logger interface {
	Debug(messageArgs ...any)
	Info(messageArgs ...any)
	Warn(messageArgs ...any)
	Error(messageArgs ...any)
	NewError(messageArgs ...any) error
}

func NewLogger(loggerName string) Logger {
	return &customLogger{loggerName}
}

func (logger *customLogger) Debug(messageArgs ...any) {
	logMessage(logger.loggerName, DEBUG, messageArgs...)
}
func (logger *customLogger) Info(messageArgs ...any) {
	logMessage(logger.loggerName, INFO, messageArgs...)
}
func (logger *customLogger) Warn(messageArgs ...any) {
	logMessage(logger.loggerName, WARN, messageArgs...)
}
func (logger *customLogger) Error(messageArgs ...any) {
	logMessage(logger.loggerName, ERROR, messageArgs...)
}
func (logger *customLogger) NewError(messageArgs ...any) error {
	msg := argsToMessage(messageArgs...)
	logMessage(logger.loggerName, ERROR, msg)
	return errors.New(msg)
}
