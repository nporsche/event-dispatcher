package dispatcher

import (
	"log"
)

type Logger interface {
	Debug(s ...interface{})
	Info(s ...interface{})
	Warning(s ...interface{})
	Error(s ...interface{})
	Fatal(s ...interface{})
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	log.Println(v)
}

func (l *DefaultLogger) Info(v ...interface{}) {
	log.Println(v)
}

func (l *DefaultLogger) Warning(v ...interface{}) {
	log.Println(v)
}

func (l *DefaultLogger) Error(v ...interface{}) {
	log.Println(v)
}

func (l *DefaultLogger) Fatal(v ...interface{}) {
	log.Fatalln(v)
}
