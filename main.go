package socketiominiservice

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zishang520/socket.io/servers/socket/v3"
	iolog "github.com/zishang520/socket.io/v3/pkg/log"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

// SocketIOService represents a Socket.IO service with generic type support
type SocketIOService[T SocketIOConstraint] struct {
	SocketIO            *socket.Server                                            // Socket.IO server instance
	middleWareFunc      socket.NamespaceMiddleware                                // Custom middleware function
	tokenValidateFunc   func(context.Context, string) (ClientInterface[T], error) // Token validation function
	userConnections     map[T][]*socket.Socket                                    // Map of user connections
	socketIOConnections map[socket.SocketId]T                                     // Map of socket connections to user IDs
	socketIOEvents      []func(context.Context) (string, func(...any))            // Registered Socket.IO events
	debug               bool                                                      // Debug mode flag
	mu                  sync.RWMutex                                              // Mutex for thread safety
}

// NewSocketIOService creates a new SocketIOService instance with default configuration
func NewSocketIOService[T SocketIOConstraint]() *SocketIOService[T] {
	socketIO := &SocketIOService[T]{
		userConnections:     make(map[T][]*socket.Socket), // Initialize user connections map
		socketIOConnections: make(map[socket.SocketId]T),  // Initialize socket connections map
		debug:               true,                         // Enable debug mode by default
	}
	opts := socket.DefaultServerOptions()

	// Configure server options
	opts.SetPingTimeout(20 * time.Second)  // Set ping timeout to 20 seconds
	opts.SetPingInterval(25 * time.Second) // Set ping interval to 25 seconds
	opts.SetMaxHttpBufferSize(1e6)         // Set max HTTP buffer size to 1MB
	opts.SetCors(&types.Cors{              // Configure CORS
		Origin:      "*",
		Credentials: true,
	})

	iolog.DEBUG = socketIO.debug

	socketIO.SocketIO = socket.NewServer(nil, opts) // Create new Socket.IO server

	return socketIO
}

// Run starts the Socket.IO service and registers middleware and events
func (s *SocketIOService[T]) Run() error {
	// Use custom middleware if provided, otherwise use default middleware
	if s.middleWareFunc == nil {
		s.SocketIO.Use(s.defaultMiddleware())
	} else {
		s.SocketIO.Use(s.middleWareFunc)
	}

	// Check if any events are registered
	if len(s.socketIOEvents) == 0 {
		return errors.New("events is empty")
	}

	// Register connection handler
	err := s.onConnect()
	if err != nil {
		return err
	}

	return nil
}

// defaultMiddleware provides default authentication and context setup for Socket.IO connections
func (s *SocketIOService[T]) defaultMiddleware() socket.NamespaceMiddleware {
	return func(client *socket.Socket, next func(*socket.ExtendedError)) {
		// Create a new context with request ID
		ctx := SetRequestIdToOrigCtx(context.Background())

		// Extract token from handshake query parameters
		token, ok := client.Handshake().Query.Query()["token"]
		if !ok {
			next(socket.NewExtendedError("error", "token is required"))
			return
		}

		// Check token format
		if !strings.HasPrefix(token[0], "Bearer ") {
			next(socket.NewExtendedError("error", "token must start with Bearer"))
			return
		}

		// Validate token using custom validation function if provided
		if s.tokenValidateFunc != nil {
			clientInfo, err := s.tokenValidateFunc(ctx, token[0])
			if err != nil {
				next(socket.NewExtendedError("error", err.Error()))
				return
			}
			if clientInfo == nil {
				next(socket.NewExtendedError("error", "client info not found"))
				return
			}
			ctx = SetClientToCtx(ctx, clientInfo)
		} else {
			// When using default token validation, T should be string
			// Convert token string to type T
			clientId, ok := any(token[0]).(T)
			if !ok {
				next(socket.NewExtendedError("error", "invalid token type"))
				return
			}
			clientInfo := &defaultToken[T]{token: clientId}
			ctx = SetClientToCtx(ctx, clientInfo)
		}

		// Store context in client data
		client.SetData(ctx)

		// Continue with next middleware
		next(nil)
	}
}

// Middleware sets a custom middleware function for the Socket.IO service
func (s *SocketIOService[T]) Middleware(f socket.NamespaceMiddleware) {
	s.middleWareFunc = f
}

// onConnect handles new Socket.IO client connections
func (s *SocketIOService[T]) onConnect() error {
	s.SocketIO.On("connection", func(clients ...any) {
		// Get client from connection arguments
		client := clients[0].(*socket.Socket)

		// Extract context from client data
		ctx, ok := client.Data().(context.Context)
		if !ok {
			return
		}

		// Get client information from context
		clientInfo := GetClientFromContext[T](ctx)
		if clientInfo == nil {
			return
		}

		// Thread-safe access to connections map
		s.mu.Lock()
		defer s.mu.Unlock()

		// Add user to the list of connected users
		if connections, ok := s.userConnections[clientInfo.GetId()]; ok {
			s.userConnections[clientInfo.GetId()] = append(connections, client)
		} else {
			s.userConnections[clientInfo.GetId()] = []*socket.Socket{client}
		}
		s.socketIOConnections[client.Id()] = clientInfo.GetId()

		// Register events for this client
		for _, fn := range s.socketIOEvents {
			client.On(fn(ctx))
		}

		// Register disconnect handler
		client.On(s.onDisconnect(ctx, client))
	})

	return nil
}

// onDisconnect handles Socket.IO client disconnections
func (s *SocketIOService[T]) onDisconnect(_ context.Context, client *socket.Socket) (string, func(...any)) {
	return "disconnect", func(data ...any) {
		// Thread-safe access to connections map
		s.mu.Lock()
		defer s.mu.Unlock()

		// Check if client exists in connections map
		if userId, ok := s.socketIOConnections[client.Id()]; ok {
			// Remove client from user connections
			if len(s.userConnections[userId]) == 1 {
				// If this is the last connection, remove the user entirely
				delete(s.userConnections, userId)
			} else {
				// Otherwise, remove just this connection from the user's connections
				if connections, ok := s.userConnections[userId]; ok {
					for i, connection := range connections {
						if connection == client {
							s.userConnections[userId] = append(connections[:i], connections[i+1:]...)
						}
					}
				}
			}
			// Remove client from socket connections map
			delete(s.socketIOConnections, client.Id())
		}
	}
}

// AddEvent registers a new event handler for the Socket.IO service
func (s *SocketIOService[T]) AddEvent(f func(context.Context) (string, func(...any))) {
	s.socketIOEvents = append(s.socketIOEvents, f)
}

// SendTo sends an event to a specific user by their ID
func (s *SocketIOService[T]) SendTo(userId T, ev string, args ...any) (bool, error) {
	// Check if user has any active connections
	if clients, ok := s.userConnections[userId]; ok {
		// Send event to all user's connections
		for _, client := range clients {
			if err := client.Emit(ev, args...); err != nil {
				return false, errors.New("sender error") // Error during sending
			}
		}
		return true, nil // Sending occurred
	}
	return false, nil // No sending occurred
}

// Handle returns the HTTP handler for the Socket.IO service
func (s *SocketIOService[T]) Handle() http.Handler {
	return s.SocketIO.ServeHandler(nil)
}

// Shutdown gracefully closes all Socket.IO connections
func (s *SocketIOService[T]) Shutdown(ctx context.Context) error {
	s.SocketIO.Close(func(err error) {})
	return nil
}
