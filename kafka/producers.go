package kafka

import (
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
)

type Producer struct {
	producer       sarama.SyncProducer
	topic          string
	schemaRegistry *srclient.SchemaRegistryClient
	codec          *goavro.Codec
	schemaID       int
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

	log.Println("Producer created with schema ID:", schema.ID())
	return &Producer{
		producer:       producer,
		topic:          topic,
		schemaRegistry: schemaRegistryClient,
		codec:          codec,
		schemaID:       schema.ID(),
	}, nil
}

func (p *Producer) SendMessage(message interface{}) error {

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return err
	}

	native, _, err := p.codec.NativeFromTextual(jsonBytes)
	if err != nil {
		log.Println("Failed to encode message:", err)
		return err
	}
	valueBytes, err := p.codec.BinaryFromNative(nil, native)
	if err != nil {
		log.Println("Failed to encode message:", err)
		return err
	}

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(p.schemaID))

	var messageValue []byte

	messageValue = append(messageValue, byte(0))
	messageValue = append(messageValue, schemaIDBytes...)
	messageValue = append(messageValue, valueBytes...)

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(messageValue),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
