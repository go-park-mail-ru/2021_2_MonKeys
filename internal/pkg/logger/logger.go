package logger

import (
	"dripapp/configs"
	"log"
	"os"
	"time"
)

type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Err   *log.Logger
}

var (
	DripLogger = Logger{
		Debug: log.New(os.Stderr, "WARNING\t", log.Ldate|log.Ltime),
		Info:  log.New(os.Stderr, "DEBUG\t", log.Ldate|log.Ltime),
		Warn:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		Err:   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
)

func (l *Logger) DebugLogging(message string) {
	if configs.LogLevel == configs.DEBUG {
		l.Debug.Printf("MESSAGE %s", message)
	}
}

func (l *Logger) InfoLogging(method string, remoteAddr string, url string, time time.Duration) {
	if configs.LogLevel <= configs.INFO {
		l.Info.Printf("[%s] %s, %s %s", method, remoteAddr, url, time)
	}
}

func (l *Logger) WarnLogging(code int, message string) {
	if configs.LogLevel <= configs.WARNING {
		l.Warn.Printf("CODE %d MESSAGE %s", code, message)
	}
}

func (l *Logger) ErrorLogging(code int, message string) {
	l.Err.Printf("CODE %d MESSAGE %s", code, message)
}
