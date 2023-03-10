### Microservices in Go

This is a sample project for a microservice architecture. Currently, the following services have been partially
implemented:

* Front-End service displays web pages to test all services.
* Broker service is an optional single point of entry into the microservice cluster.
* Authentication service is a simple authentication service, which uses Chi as the router and Postgres as the database.

The following services will be added later:

* Logging service: Uses a MongoDB database.
* Listener service: Receives messages from RabbitMQ and acts upon them.
* Mail service: Takes a JSON payload, converts into a formatted email, and sends it out.

#### Table of Contents

<div class="toc">
  <ul>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#project-structure">Project Structure</a></li>
    <li><a href="#usage">Usage</a>
      <ul>
        <li><a href="#front-end">Front-end</a></li>
        <li><a href="#broker-service">Broker service</a></li>
        <li><a href="#authentication-service">Authentication Service</a></li>
      </ul>
    </li>
  </ul>
</div>

#### Getting Started

Prerequisites:

* Go 1.18
* Postgres
* Docker

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

The App should now be accessible at `http://localhost`.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

#### Project Structure

```
.
├── authentication-service
│ ├── cmd
│ │ └── api
│ │ ├── handlers.go
│ │ ├── helpers.go
│ │ ├── main.go
│ │ └── routes.go
│ ├── data
│ │ └── models.go
│ ├── authentication-service.dockerfile
│ └── go.mod
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

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

#### Usage

##### Front-end

This project is a front-end service that interacts with microservices in Go. It provides a simple web page with two
buttons, one for testing the broker service and one for testing the authentication service.

To use the Front End service, navigate to http://localhost in a web browser.

This is a microservices testing application that tests two endpoints: Test Broker and Test Auth.

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Broker Service

The broker service serves as a proxy between clients and various backend services. It is responsible for authenticating
clients and routing requests to the appropriate backend service.

**Endpoints**

The broker service provides two endpoints:

* `/`
* `/handle`

`GET /` is a simple endpoint that returns a response indicating that the broker service is up and running.

Response:

```bash
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

<p align="right">(<a href="#table-of-contents">back to the Table of content</a>)</p>

##### Authentication Service

This is a simple authentication service built in Go, which uses Chi as the router and Postgres as the database.

**Endpoints**

The following endpoints are available:

`POST /authenticate` - authenticate a user.

The request payload should contain the following fields: email and password. If the credentials are valid, the endpoint
will return a JSON response containing the user object and a success message. If the credentials are invalid, it will
return an error message.

Example request:

```bash
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

```bash
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

```bash
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

