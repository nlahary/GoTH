This project is a simple implementation of a website backend in golang, using the following stack:
- `http/net` package for the routing of the website
- golang templates for the frontend
- sqlite3 for the database
- redis for the cookie storage, to keep track of the user cart

Additionally, the project uses Kafka for the logging of the code execution and the http requests, using respectively the `logs` and `httplogs` topics.
The messages are serialized using Avro and consumed by an ElasticSearch sink, which is used to visualize the logs in Kibana.

For the Kafka implementation, the project uses the following libraries:
- `github.com/IBM/sarama` for the Kafka producers 
- `github.com/riferrei/srclient` for the Schema Registry client
- `github.com/linkedin/goavro/v2` for the Avro serialization 

The project is containerized using Docker and orchestrated using Docker Compose. The services are defined in the `docker-compose.yml` file, and the project can be started using the `start.sh` script.

The script will start all the services and push the ElasticSearch sink configuration to the Kafka Connect service.

The Kafka services can be monitored using the KafkaUI service, which is available at `localhost:8080`.