#!/bin/bash

SINK_CONNECTOR_PATH=$PWD/kafka/sink/elasticsearch-sink.json

if [ ! -f "$SINK_CONNECTOR_PATH" ]; then
  echo "File not found!"
  exit 1
fi

# Start the services
docker-compose up -d 

  response=$(curl -s -X GET http://localhost:8083/connectors/elasticsearch-sink/status)
  if [[ $response == *"RUNNING"* ]]; then
    break
  fi
  echo "waiting for kafka connector to be up and running"
  sleep 10
done

curl -s -X POST -H "Content-Type: application/json" --data @"$SINK_CONNECTOR_PATH" http://localhost:8083/connectors

response=$(curl -s -X GET http://localhost:8083/connectors/elasticsearch-sink/status)

while [[ $response != *"RUNNING"* ]]; do
  sleep 5
  response=$(curl -s -X GET http://localhost:8083/connectors/elasticsearch-sink/status)
done

echo "Connector is running!"

air