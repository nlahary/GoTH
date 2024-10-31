package models

import (
	"fmt"
	"log"
	"time"
)

// BasicLogger is a logger that inherits from DefaultLogger
// used to log code execution messages.
// It overrides the Print, Println, Fatal and Error methods
// from the log package to send messages to Kafka.
type BasicLogger struct {
	*DefaultLogger
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

func (l *BasicLogger) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "INFO",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Print(v...) // Pour un retour visuel immédiat, optionnel
	l.SendToKafka(l.Map(logMsg))
}

func (l *BasicLogger) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	logMsg := LogMessage{
		Level:     "INFO",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Println(v...) // Pour un retour visuel immédiat, optionnel
	l.SendToKafka(l.Map(logMsg))
}

func (l *BasicLogger) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "FATAL",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	l.SendToKafka(l.Map(logMsg))
	log.Fatal(v...) // Arrête le programme

}

func (l *BasicLogger) Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	logMsg := LogMessage{
		Level:     "ERROR",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	log.Print(v...) // Pour un retour visuel immédiat, optionnel
	l.SendToKafka(l.Map(logMsg))
}
