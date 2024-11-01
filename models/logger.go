package models

import (
	"fmt"
	"log"
	"time"

	kafka "github.com/nlahary/website/kafka"
)

type CodeExecLogger struct {
	kafka.Producer
}

const BasicLogSchema = `{
	"type": "record",
	"name": "LogMessage",
	"fields": [
		{"name": "level", "type": "string"},
		{"name": "message", "type": "string"},
		{"name": "timestamp", "type": "string"}
	]
}`

type LogMessage struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func (l *CodeExecLogger) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "INFO",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Print(v...) // Pour un retour visuel immédiat, optionnel
	MappedLogMsg := l.Producer.Map(logMsg)
	l.Producer.SendMessage(MappedLogMsg)
}

func (l *CodeExecLogger) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	logMsg := LogMessage{
		Level:     "INFO",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Println(v...) // Pour un retour visuel immédiat, optionnel
	MappedLogMsg := l.Producer.Map(logMsg)
	l.Producer.SendMessage(MappedLogMsg)
}

func (l *CodeExecLogger) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "FATAL",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	MappedLogMsg := l.Producer.Map(logMsg)
	l.Producer.SendMessage(MappedLogMsg)
	log.Fatal(v...) // Arrête le programme

}

func (l *CodeExecLogger) Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "ERROR",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Print(v...) // Pour un retour visuel immédiat, optionnel
	MappedLogMsg := l.Producer.Map(logMsg)
	l.Producer.SendMessage(MappedLogMsg)
}
