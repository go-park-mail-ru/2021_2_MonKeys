package logger

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
	Warn *log.Logger
}

var (
	InfoLog    = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog   = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	WatningLog = log.New(os.Stderr, "WARNING\t", log.Ldate|log.Ltime)

	DripLogger = Logger{
		Info: InfoLog,
		Err:  ErrorLog,
		Warn: WatningLog,
	}
)

func (l *Logger) InfoLogging(method string, remoteAddr string, url string, time time.Duration) {
	l.Info.Printf("[%s] %s, %s %s", method, remoteAddr, url, time)
}

func (l *Logger) ErrorLogging(code int, message string) {
	l.Err.Printf("CODE %d MESSAGE%s", code, message)
}

func (l *Logger) WarnLogging(code int, message string) {
	l.Warn.Printf("CODE %d MESSAGE%s", code, message)
}
