Certainly! Here is the `README.md` file formatted for easy copying:

```markdown
# rstream-consumer

`rstream-consumer` is a Go-based application designed to consume and process messages from a Redis stream. It includes various internal packages for configuration, consumer logic, storage, shutdown handling, and monitoring.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Testing](#testing)
- [Docker](#docker)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Message Consumption**: Efficiently consume messages from Redis streams.
- **Monitoring**: Monitor message processing with detailed logs.
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
go mod download
```

## Usage

### Running the Application

```sh
go run main.go consumer
```

### Building the Application

```sh
go build -o rstream-consumer main.go
```

## Configuration

Configuration can be provided via environment variables or a configuration file. Refer to the `internal/config` package for more details.

## Testing

### Running Tests

```sh
go test -cover ./...
```

### Coverage Report

To generate a coverage report:

```sh
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out -o coverage.html
```

## Docker

### Building the Docker Image

```sh
docker build -t rstream-consumer -f .deployment/docker/Dockerfile .
```

### Running the Docker Container

```sh
docker run --rm -e REDIS_URL=redis://localhost:6379 -p 8080:8080 rstream-consumer
```


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```