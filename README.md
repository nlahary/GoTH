<p align="center">
  <img src="images/golang.png" alt="golang Logo" width="200", style="margin-right: 20px">
  <img src="images/kafka.png" alt="kafka Logo" width="165", style="margin-left: 20px">
</p>

This project is a simple implementation of a web server in golang, using the following stack:
- `http/net` package for the routing of the API
- golang templates + HTMX for the frontend
- sqlite3 for the database
- redis for the cookie storage, to keep track of the user cart

The API has the following endpoints:
1. `/`
    - GET (default): Returns the home page with the list of contacts and CRUD operations
    
2.  `/contacts`
    - POST: Add a new contact to the database with a form submission 

3.  `/contacts/{contact_id}`
    - GET: Returns the contact details (name and email)
    - DELETE: Delete a contact from the database
    - PUT: Update a contact in the database with a form submission

4.  `/contacts/{contact_id}/edit`
    - GET: Returns the contact details in a form for editing

5. `/Products` - Products page:
    - GET: Returns all the products in the database and displays them on the page

6. `/cart/{product_id}`
    - POST: Add a product to the user cart

Additionally, the project uses a Kafka broker for the logging of the code execution and the http requests, using respectively the `logs` and `httplogs` topics.

The messages are de/serialized using Avro and consumed by an ElasticSearch sink. After a quick configuration, you can see the logs using the Kibana service, available at `localhost:5601`.

For the Kafka implementation, the project uses the following libraries:
- `github.com/IBM/sarama` for the Kafka producers 
- `github.com/riferrei/srclient` for the Schema Registry client
- `github.com/linkedin/goavro/v2` for the Avro serialization 

The project is containerized using Docker and orchestrated using Docker Compose. The services are defined in the `docker-compose.yml` file, and the project can be started using the `start.sh` script.

The script will start all the services and push the ElasticSearch sink configuration to the Kafka Connect service.

The Kafka services can be monitored using the KafkaUI service, which is available at `localhost:8080`.

To run the project:
1. Clone the repository
2. Run the `start.sh` script
3. Wait for the services to be up and running (It will be log in the terminal when ready)
4. Open the website at `localhost:42069`

The script will wait for the services to be up and running, post the ElasticSearch sink configuration, and start the website at `localhost:42069` using `air` package for hot reloading.

The kafka broker and schema registry might crash for some reason, restart them and re-run the script in case.


