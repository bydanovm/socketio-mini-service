# SocketIO Mini Service

[![ru](https://img.shields.io/badge/lang-ru-yellow.svg)](https://github.com/bydanovm/socketio-mini-service/blob/main/README.RU.md)

SocketIO Mini Service is a minimalist service for working with Socket.IO, created to simplify the development of real-time applications. This service provides a generic interface for managing Socket.IO connections with support for custom authorization and event handling.

The foundation of the service is built on [Socket.IO by zishang520](https://github.com/zishang520/socket.io).

## Key Features

- Management of Socket.IO connections with support for generic types
- Built-in token-based authorization system
- Support for custom middleware
- Thread-safe connection management
- Flexible event handling system
- CORS support
- Configurable timeout and buffering parameters

## Installation

```bash
go get github.com/bydanovm/socketio-mini-service
```

## Usage

Example of creating a basic service:

```go
package main

import (
    "context"
    "net/http"
    "log"
    
    "github.com/bydanovm/socketio-mini-service"
)

type User struct {
    ID   string
    Name string
}

func (u *User) GetId() string {
    return u.ID
}

func main() {
    // Create a new service
    service := socketiominiservice.NewSocketIOService[string]()
    
    // Set token validation function
    service.TokenValidate(func(ctx context.Context, token string) (socketiominiservice.ClientInterface[string], error) {
        // Validate token and return user info
        // This is a simplified example
        if token == "Bearer valid-token" {
            return &User{ID: "user1", Name: "John Doe"}, nil
        }
        return nil, errors.New("invalid token")
    })
    
    // Add event handler
    service.AddEvent("message", func(ctx context.Context, data socketiominiservice.SocketIOData[string]) {
        // Handle message event
        log.Printf("Received message from user %s: %s", data.Client.GetId(), string(data.Payload))
    })
    
    // Run the service
    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
    
    // Start HTTP server
    http.Handle("/socket.io/", service.Handle())
    log.Fatal(http.ListenAndServe(":3000", nil))
}
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Main Developer

**bydanovm** - the main developer of the minimalist service.

## Contributors

The project is open for contributions from the developer community.
