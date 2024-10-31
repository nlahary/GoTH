package models

import (
	"log"
	"reflect"
	"strings"

	"github.com/linkedin/goavro/v2"
	kafka "github.com/nlahary/website/kafka"
)

// DefaultLogger is a logger that sends messages to Kafka
// using Avro serialization
type DefaultLogger struct {
	producer *kafka.Producer
	codec    *goavro.Codec
}

func NewLogger(producer *kafka.Producer, schema string) *DefaultLogger {

	codec, err := goavro.NewCodec(string(schema))
	if err != nil {
		log.Fatalf("Erreur de création du codec Avro: %v", err)
	}

	return &DefaultLogger{
		producer: producer,
		codec:    codec,
	}
}

// Given a struct supposed to be a message to send to Kafka,
// returns a map representation of it for serialization.
// The map keys are the struct field names, with respect to
// the JSON tags.
func (l *DefaultLogger) Map(message interface{}) map[string]interface{} {
	MappedMessage := make(map[string]interface{})
	if message == nil {
		log.Fatalf("Message vide")
		return MappedMessage
	}
	value := reflect.ValueOf(message)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		log.Fatalf("Message invalide")
		return MappedMessage
	}
	typeOfStruct := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typeOfStruct.Field(i)
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			fieldName = jsonTag
		} else {
			fieldName = strings.ToLower(fieldName)
		}
		MappedMessage[fieldName] = value.Field(i).Interface()
	}
	return MappedMessage
}

func (l *DefaultLogger) SendToKafka(message map[string]interface{}) {

	avroMsg, err := l.codec.TextualFromNative(nil, message)
	if err != nil {
		println(message)
		log.Fatalf("Erreur de sérialisation Avro: %v", err)
	}
	if err := l.producer.SendMessage(string(avroMsg)); err != nil {
		log.Fatalf("Erreur d'envoi du message: %v", err)
	}
}
