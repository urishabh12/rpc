package erros

const (
	connectionClosed = "connection has been closed"
)

type ConnClosedError struct{}

func (c ConnClosedError) Error() string {
	return connectionClosed
}

func NewConnClosedError() ConnClosedError {
	return ConnClosedError{}
}

func IsConnClosedError(err error) bool {
	return err.Error() == connectionClosed
}
