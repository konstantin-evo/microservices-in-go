![Go Language](https://img.shields.io/badge/Go-1.18-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10.9-blue.svg)
![MongoDB](https://img.shields.io/badge/MongoDB-4.4.9-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13.5-blue.svg)
![MailHog](https://img.shields.io/badge/MailHog-1.0.1-blue.svg)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.9-blue.svg)

### Microservices in Go

<p align="center">
  <img src="project/src/img/demo-header.png" width="70%" alt="demo header">
</p>

<!-- CONTACTS -->
<div>
    <p align="center">
        <br />
        <a href="#about-project">About Project</a>
        ·
        <a href=mailto:"konstantin.priluchnyi@gmail.com?subject=report%20Bug">Report Bug</a>
        ·
        <a href=mailto:"konstantin.priluchnyi@gmail.com?subject=Request%20Feature">Request Feature</a>
    </p>
    <p align="center">
        <a href=mailto:"konstantin.priluchnyi@gmail.com">
            <img src="https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white">
        </a>
        <a href="https://t.me/konstantin_evo">
            <img src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white">
        </a>
        <a href="https://www.linkedin.com/in/konstantin-evo/">
            <img
                src="https://img.shields.io/badge/linkedin-%230077B5.svg?style=for-the-badge&logo=linkedin&logoColor=white">
        </a>
    </p>
</div>

#### Table of Contents

<div class="toc">
  <ul>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#about-project">About project</a></li>
      <ul>
        <li><a href="#project-structure">Project Structure</a></li>
        <li><a href="#environment-variables">Environment Variables</a></li>
      </ul>
    <li><a href="#usage">Usage</a>
      <ul>
        <li><a href="#front-end">Front-end</a></li>
        <li><a href="#broker-service">Broker service</a></li>
        <li><a href="#authentication-service">Authentication Service</a></li>
        <li><a href="#logger-service">Logger Service</a></li>
        <li><a href="#mail-service">Mail Service</a></li>
        <li><a href="#listener-service">Listener Service</a></li>
      </ul>
    </li>
  </ul>
</div>

#### Getting Started

Prerequisites:

* Go 1.18
* Docker
* MongoDB Compass (only to see data in MongoDB database via GUI)

To get started, clone this repository and navigate to the project directory:

```
git clone https://github.com/konstantin-evo/microservices-in-go.git
cd microservices-in-go
```

1. Start the project using docker-compose
2. Build front-end app
3. Run app

```
cd ./project && docker-compose up -d
cd ./../front-end && env CGO_ENABLED=0 go build -o frontApp ./cmd/web
./frontApp 
```

Endpoints:

| Service     | URL                     | Credentials                                                              | Description                        |
|-------------|-------------------------|--------------------------------------------------------------------------|------------------------------------|
| Front-end   | http://localhost        | -                                                                        | The main application interface     |
| Mail server | http://localhost:8025/  | -                                                                        | MailHog UI to view captured emails |
| Database    | http://localhost:8089/  | POSTGRES_USER: postgres, POSTGRES_PASSWORD: password, POSTGRES_DB: users | Access to the users database       |
| RabbitMQ    | http://localhost:15672/ | guest:guest                                                              | RabbitMQ management console        |

**Note**: To view the data in the "logs" database (used in Logging Service), it is necessary to install MongoDB Compass.

<p>You can find more information about connection to PostgreSQL and MongoDB in the <a href="#environment-variables">Environment Variables</a> section</p>

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

#### About project

The "Microservices in Go" project is an open-source code repository that demonstrates how to build microservices using
the Go programming language. The repository consists of several modules, each of which is a separate microservice.

The project provides an example of how to build microservices using Go, and it also demonstrates various best
practices and design patterns for microservice architecture.

Currently, the following services have been partially implemented:

| Service        | Description                                                                               |
|----------------|-------------------------------------------------------------------------------------------|
| Front-End      | Displays web page to test all services.                                                   |
| Broker         | Routes requests to the appropriate backend service.                                       |
| Authentication | A simple authentication service that uses Chi as the router and Postgres as the database. |
| Logging        | Allows clients to write log entries to a MongoDB database.                                |
| Mail           | Takes a JSON payload, converts into a formatted email, and sends it out.                  |
| Listener       | Receives messages from RabbitMQ and acts upon them.                                       |

```mermaid
graph TD
    subgraph Frontend service
        Frontend((Frontend))
    end

    subgraph Broker service
        Broker((Broker Service))
    end

    subgraph Authentication service
        Auth((Auth Service))
        Database
    end

    subgraph Logger service
        Logger((Logger Service))
        Database_NoSQL
    end

    subgraph Mail service
        Mail((Mail Service))
    end

    subgraph Listener service
        Listener((Listener Service))
    end

    subgraph Database
        PostgreSQL
    end

    subgraph Database_NoSQL
        MongoDB
    end

    Frontend -- HTTP/HTTPS --> Broker
    Broker -- HTTP/HTTPS --> Auth
    Broker -- HTTP/HTTPS --> Mail
    Broker -- AMQP --> RabbitMQ
    Broker -- RPC --> Logger
    RabbitMQ -- AMQP --> Listener
    Listener -- HTTP/HTTPS --> Logger
    Auth -- HTTP/HTTPS --> Logger
```

##### Project Structure

This project represents a multi-service application consisting of a front-end, broker-service, authentication-service,
logging-service, and mail-service. Each service has its own Dockerfile and Go module configuration. The directory
structure is as follows:

```
.
├── authentication-service
│   ├── authentication-service.dockerfile
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── main.go
│   │       └── routes.go
│   ├── data
│   │   └── models.go
│   └── go.mod
├── broker-service
│   ├── broker-service.dockerfile
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── main.go
│   │       └── routes.go
│   ├── data
│   │   └── models.go
│   ├── event
│   │   ├── consumer.go
│   │   ├── data
│   │   │   └── models.go
│   │   ├── emitter.go
│   │   ├── event.go
│   │   └── logger.go
│   └── go.mod
├── front-end
│   ├── cmd
│   │   └── web
│   │       ├── main.go
│   │       └── templates
│   │           ├── base.layout.gohtml
│   │           ├── footer.partial.gohtml
│   │           ├── header.partial.gohtml
│   │           └── test.page.gohtml
│   └── go.mod
├── listener-service
│   ├── event
│   │   ├── consumer.go
│   │   ├── event.go
│   │   └── logger.go
│   ├── go.mod
│   ├── go.sum
│   ├── listener-service.dockerfile
│   └── main.go
├── logging-service
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── main.go
│   │       ├── rpc.go
│   │       └── routes.go
│   ├── data
│   │   └── models.go
│   ├── go.mod
│   └── logger-service.dockerfile
├── mail-service
│   ├── cmd
│   │   └── api
│   │       ├── handlers.go
│   │       ├── helpers.go
│   │       ├── mailer.go
│   │       ├── main.go
│   │       └── routes.go
│   ├── go.mod
│   ├── mail-service.dockerfile
│   └── templates
│       ├── mail.html.gohtml
│       └── mail.plain.gohtml
└── project
  ├── Makefile
  └── docker-compose.yml
```

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Environment Variables

**Authentication Service**

The authentication service uses the `DSN` environment variable that contains the connection string for the PostgreSQL
database, including host, port, user, password, database name, SSL mode, timezone, and connection timeout.

Example:

```
host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5
```

**Mailer Service**

The mailer service uses the following environment variables:

| Variable        | Description                                    | Example                    |
|-----------------|------------------------------------------------|----------------------------|
| MAIL_DOMAIN     | The domain name of the email server.           | localhost                  |
| MAIL_HOST       | The hostname of the email server.              | mailhog                    |
| MAIL_PORT       | The port number of the email server.           | 1025                       |
| MAIL_ENCRYPTION | The encryption method for email.               | none                       |
| MAIL_USERNAME   | The username for the email server.             | ""                         |
| MAIL_PASSWORD   | The password for the email server.             | ""                         |
| FROM_NAME       | The display name for the sender of the emails. | "Konstantin Evo"           |
| FROM_ADDRESS    | The email address of the sender.               | konstantin.evo@example.com |

**PostgreSQL**

The PostgreSQL service uses the following environment variables:

| Variable          | Description                               | Example  |
|-------------------|-------------------------------------------|----------|
| POSTGRES_USER     | The username for the PostgreSQL database. | postgres |
| POSTGRES_PASSWORD | The password for the PostgreSQL database. | password |
| POSTGRES_DB       | The name of the PostgreSQL database.      | users    |

**MongoDB**

The MongoDB service uses the following environment variables:

| Variable                   | Description                                 | Example  |
|----------------------------|---------------------------------------------|----------|
| MONGO_INITDB_DATABASE      | The name of the MongoDB database.           | logs     |
| MONGO_INITDB_ROOT_USERNAME | The root username for the MongoDB database. | admin    |
| MONGO_INITDB_ROOT_PASSWORD | The root password for the MongoDB database. | password |

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

#### Usage

##### Front-end

This project is a front-end service that interacts with microservices in Go.

The application interface consists of buttons to test each microservice and fields to input authentication details.

To test a microservice, click on the corresponding button. The application will send a request to the microservice and
display the response in the "Received" section. The request payload will be displayed in the "Sent" section.

If an error occurs, the error message will be displayed in the "Logs" section.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

![demo.png](project%2Fsrc%2Fimg%2Fdemo.png)

##### Broker Service

The broker service serves as a proxy between clients and various backend services. It is responsible for routing
requests to the appropriate backend service.

The broker service utilizes RabbitMQ to handle messaging with logger service and also supports RPC communication for
specific interactions with the Logger Service.

**Endpoints**

The broker service provides two endpoints:

* `/`
* `/handle`

`GET /` is a simple endpoint that returns a response indicating that the broker service is up and running.

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "error": false,
  "message": "Hit the broker"
}
```

`POST /handle`

This endpoint is used to handle requests from clients. The broker service expects a JSON payload with an action field
indicating the type of action to perform, and an optional auth field containing authentication information.

The request and response depend on the `action` field. The example for authentication:

```json
{
  "action": "auth",
  "auth": {
    "email": "user@example.com",
    "password": "password123"
  }
}
```

If the action field is not recognized, the response will contain an error message:

```bash
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "error": true,
  "message": "unknown action"
}
```

For the Logger Service, the Broker Service supports both RabbitMQ messaging and RPC communication. The Broker Service
uses the logEvent function for RabbitMQ messaging and the logItemViaRPC function for RPC communication.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Authentication Service

This is a simple authentication service built in Go, which uses Chi as the router and Postgres as the database.

The "users" database containing user data can be accessed at http://localhost:8089/

To connect to the database, you need to provide the following information:

* **server**: postgres
* **username**: postgres
* **password**: password
* **database**: users

![demo-postresql-adminer.png](project%2Fsrc%2Fimg%2Fdemo-postresql-adminer.png)

**Endpoints**

The following endpoints are available:

`POST /authenticate` - authenticate a user.

The request payload should contain the following fields: email and password. If the credentials are valid, the endpoint
will return a JSON response containing the user object and a success message. If the credentials are invalid, it will
return an error message.

Example request:

```
POST /authenticate HTTP/1.1
Host: localhost:80
Content-Type: application/json
```

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

Example response:

```
HTTP/1.1 202 Accepted
Content-Type: application/json
```

```json
{
  "error": false,
  "message": "Logged in user john@example.com",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "active": 1,
    "created_at": "2022-10-15T14:35:00Z",
    "updated_at": "2022-10-15T14:35:00Z"
  }
}
```

If the email and password are invalid, the server will respond with an error message:

```
HTTP/1.1 400 Bad Request
Content-Type: application/json
```

```json
{
  "error": true,
  "message": "invalid credentials"
}
```

**Structure**

The code is structured as follows:

* `cmd/api/main.go` - the main entry point for the application.
* `cmd/api/routes.go` - the routing and middleware configuration for the application.
* `cmd/api/handlers.go` - the request handlers for the endpoints.
* `cmd/api/helpers.go` - some helper functions for parsing JSON, writing JSON responses, and handling errors.
* `data/models.go` - the database models for the application.
* `authentication-service.dockerfile` - the Dockerfile for the application.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Logger Service

The Logger Service is a versatile logging system that allows clients to write log entries to a MongoDB database. It
supports two primary modes of communication: REST (HTTP) and RPC.

Different services can use the Logger Service in different ways:

1. The Auth Service communicates with the Logger Service via HTTP to send a POST request when a user is successfully
   authenticated.
2. On the other hand, the Broker Service utilizes RPC communication to interact with the Logger Service remotely, by
   calling the LogInfo method of the RPCServer struct.

In summary, the Logger Service provides flexibility in communication methods, enabling various services to interact with
it seamlessly, whether it be through REST (HTTP) or RPC.

**Endpoints**

The `POST /log` endpoint is used to write log entries to the MongoDB database. It expects a JSON payload that contains
the
name and data fields.

Request Example:

```
POST /log HTTP/1.1
Content-Type: application/json
```

```json
{
  "name": "info",
  "data": "This is an info message."
}
```

Response Example:

```
HTTP/1.1 202 Accepted
Content-Type: application/json
```

```json
{
  "error": false,
  "message": "The message is successfully logged in the service.",
  "data": null
}
```

**Structure**

The code is structured as follows:

* `cmd/api/main.go` - the main entry point for the application.
* `cmd/api/routes.go` - the routing and middleware configuration for the application.
* `cmd/api/rpc.go` - the RPC server implementation and related functions.
* `cmd/api/handlers.go` - the request handlers for the endpoints.
* `cmd/api/helpers.go` - some helper functions for parsing JSON, writing JSON responses, and handling errors.
* `data/models.go` - the database models for the application.
* `logger-service.dockerfile` - the Dockerfile for the application.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Mail Service

This is a simple mail service API, which allows you to send emails using the SMTP protocol.

Features:

1. Send emails with attachments
2. Supports HTML and plain text email templates
3. Environmental variable configuration for mail settings
4. Containerization with Docker

Demo:

![demo-mail-service.gif](project%2Fsrc%2Fimg%2Fdemo-mail-service.gif)

**Endpoints**

The following endpoints are available:

1. `POST /send`: Send an email
2. `GET /ping`: Health check endpoint to verify the service is running

To send an email, make a `POST` request to `/send` with the following JSON payload:

```json
{
  "from": "example@example.com",
  "to": "recipient@example.com",
  "subject": "Hello, World!",
  "message": "This is a test email."
}
```

Response Payload

```json
{
  "error": false,
  "message": "The message successfully sent to recipient@example.com"
}
```

To verify the service is running, make a `GET` request to `/ping`. If the service is running correctly, it will return a
200 OK status.

The code is organized into the following files:

**Structure**

The code is structured as follows:

* `handlers.go`: Contains the HTTP handler functions for the API endpoints
* `routes.go`: Configures the API routes and CORS settings
* `mailer.go`: Contains the email sending logic and email template rendering
* `helpers.go`: Contains helper functions for JSON input/output and error handling
* `main.go`: Initializes the mail service configuration and starts the HTTP server
* `go.mod` and `go.sum`: Go module dependency management files
* `mail-service.dockerfile`: Dockerfile to containerize the mail service
* `templates/mail.plain.gohtml`: Plain text email template
* `templates/mail.html.gohtml`: HTML email template

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Listener Service

This is a simple application to listen for and consume messages from RabbitMQ based on specific topics. The application
is written in Go and uses the RabbitMQ Go client library for communication with RabbitMQ.

Features:

1. Connect to a RabbitMQ server.
2. Listen for messages with specific topics.
3. Forward the consumed messages to a log service.

**Structure**

The code is structured as follows:

* `main.go`: Initializes the configuration and connects to RabbitMQ, sets up the consumer, and starts listening for
  messages.
* `event/consumer.go`: Contains the Consumer struct and methods for setting up and consuming messages from RabbitMQ.
* `event/exchange.go`: Contains the function for declaring the topic exchange.
* `event/queue.go`: Contains the function for declaring a random queue and binding it to the topic exchange.
* `event/log_event.go`: Contains the function for sending the consumed messages to the log service.
* `go.mod` and `go.sum`: Go module dependency management files.

**Configuration**

The application uses a simple configuration struct to set up the connection to RabbitMQ and the log service URL:

```
app := config{
  RabbitMqURL:   "amqp://guest:guest@rabbitmq",
  LogServiceURL: "http://logger-service/log",
  Topics:        []string{"log.INFO", "log.WARNING", "log.ERROR"},
}
```

You can modify these values according to your environment.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>
