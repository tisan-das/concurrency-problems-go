package file

type File interface {
	Open() error
	Close() error
	Read(int, int) (string, error)
	Write(int, string) error

	// TODO: Use metadata structure rather than string
	GetMetadata() string
}
