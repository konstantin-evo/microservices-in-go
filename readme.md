### Microservices in Go

This is a sample project for a microservice architecture. Currently, the following services have been partially
implemented:

* Broker service: An optional single point of entry into the microservice cluster.
* Front End service: Displays web pages.

The following services will be added later:

* Authentication service: Uses a Postgres database.
* Logging service: Uses a MongoDB database.
* Listener service: Receives messages from RabbitMQ and acts upon them.
* Mail service: Takes a JSON payload, converts into a formatted email, and sends it out.

#### Getting Started

To get started, clone this repository and navigate to the project directory:

```
git clone https://github.com/example/microservice-project.git
cd microservice-project
```

Then, start the project using docker-compose:

```
docker-compose up -d
```

The Broker service should now be accessible at http://localhost:8080.

#### Project Structure

```
.
├── broker-service
│ ├── cmd
│ │ └── api
│ │ ├── handlers.go
│ │ ├── helpers.go
│ │ ├── main.go
│ │ └── routes.go
│ ├── go.mod
│ └── broker-service.dockerfile
├── front-end
│ ├── cmd
│ │ └── web
│ │ ├── main.go
│ │ └── templates
│ │ ├── base.layout.gohtml
│ │ ├── footer.partial.gohtml
│ │ ├── header.partial.gohtml
│ │ └── test.page.gohtml
│ ├── go.mod
└── project
└── docker-compose.yml
```

#### Usage

To use the Front End service, navigate to http://localhost:8080 in a web browser.

Currently, the Front End service only has one page: "Test microservices". This page allows you to test the Broker
service by clicking the "Test Broker" button. When clicked, the Front End service sends a POST request to the Broker
service and displays the response.

The Broker service responds with a JSON object containing an "error" flag (indicating success or failure), a "message"
field (containing a string message), and an optional "data" field (containing additional data).