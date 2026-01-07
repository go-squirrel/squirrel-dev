package terminal

type TerminalHandler interface {
	Write(data []byte) (int, error)
	Read(output []byte) (int, error)
	Resize(width, height int) error
	Close() error
}

