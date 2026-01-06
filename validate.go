package socketiominiservice

// defaultToken is a default implementation of the ClientInterface
// It simply stores and returns a token of generic type T
type defaultToken[T SocketIOConstraint] struct {
	token T
}

// GetId returns the token stored in the defaultToken
func (d *defaultToken[T]) GetId() T {
	return d.token
}
