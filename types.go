package socketiominiservice

type SocketIOData[T SocketIOConstraint] struct {
	Payload []byte
	Client  ClientInterface[T]
}
