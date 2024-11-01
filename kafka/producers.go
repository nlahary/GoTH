package kafka

import (
	"log"
	"reflect"
	"strings"

	"github.com/IBM/sarama"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

type Producer struct {
	producer       sarama.SyncProducer
	topic          string
	schemaRegistry *srclient.SchemaRegistryClient
	codec          *goavro.Codec
}

func NewProducer(brokers []string, topic string, schemaRegistryClient *srclient.SchemaRegistryClient) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	schema, err := schemaRegistryClient.GetLatestSchema(topic + "-value")
	if err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(schema.Schema())
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer:       producer,
		topic:          topic,
		schemaRegistry: schemaRegistryClient,
		codec:          codec,
	}, nil
}

// Map a struct (supposed to be a message) to a map[string]interface{}
func (p *Producer) Map(message interface{}) map[string]interface{} {
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

func (p *Producer) SendMessage(message map[string]interface{}) error {
	avroMsg, err := p.codec.BinaryFromNative(nil, message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(avroMsg),
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
