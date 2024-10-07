# rstream-consumer

`rstream-consumer` is a Go-based application designed to consume and process messages from a Redis pubsub and publishes them into a Redis Stream. It includes various internal packages for configuration, consumer logic, storage, shutdown handling, and monitoring.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Testing](#testing)
- [Docker](#docker)
- [License](#license)

## Features

- **Message Consumption**: Efficiently consume messages from Redis streams using Redis sets to make sure only one instance of the application processes a message.
- **Monitoring**: Monitor messages processing per second.
- **Graceful Shutdown**: Handle shutdown signals gracefully.
- **Configurable**: Easily configurable via environment variables or configuration files.
- **Testing**: Comprehensive test coverage.

## Installation

### Prerequisites

- Go 1.23.2 or later
- Redis

### Clone the Repository

```sh
git clone https://github.com/cuctemeh/rstream-consumer.git
cd rstream-consumer
```

### Install Dependencies

```sh
go mod tidy
```

## Usage

### Running the Application

```sh
go run main.go start
```
The application starts consuming messages from the Redis pubsub and publishes them into a Redis Stream. The application can be stopped by sending a `SIGINT` or `SIGTERM` signal.

In order to publish a predefined set of messages to the Redis pubsub you need to setup your python environment first, install all the requirements from `./scripts/producer/requirements.txt` and run the following code:
```python ./scripts/producer/server.py```

### Building the Application

```sh
go build -o rstream-consumer main.go
```

## Configuration

Configuration can be provided via environment variables or a configuration file. Refer to the `internal/config` package for more details.

## Testing

### Running Tests

```sh
go test ./...
```

### Coverage Report

To generate a coverage report:

```sh
./scripts/coverage_check.sh
```

## Docker

### Building the Docker Image
You need to have docker or podman installed on your machine to build the image. You can build the image using the following command:
```sh
docker build -t rstream-consumer-final:latest -f .deployment/docker/Dockerfile .
```
After building the image you can run multiple containers using the docker-compose.yaml with the following command:

```sh
docker-compose up -d
```
This will start 5 instances of the application with their configuration and a redis instance. 
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
