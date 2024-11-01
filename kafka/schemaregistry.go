package kafka

import (
	"fmt"

	"github.com/riferrei/srclient"
)

func RegisterSchemaIfNotExists(client *srclient.SchemaRegistryClient, topic, schemaDefinition string) error {
	subject := topic + "-value"
	_, err := client.GetLatestSchema(subject)
	if err != nil {
		_, err := client.CreateSchema(subject, schemaDefinition, srclient.Avro)
		if err != nil {
			return fmt.Errorf("failed to create schema for subject %s: %w", subject, err)
		}
	}
	return nil
}
