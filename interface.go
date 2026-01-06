package socketiominiservice

// SocketIOConstraint is a type constraint for Socket.IO client IDs
// It requires the type to be comparable
type SocketIOConstraint interface {
	comparable
}

// ClientInterface defines the interface for Socket.IO clients
type ClientInterface[T SocketIOConstraint] interface {
	// GetId returns the client's unique identifier
	GetId() T
}
