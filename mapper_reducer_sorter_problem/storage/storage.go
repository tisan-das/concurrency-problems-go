package storage

type File interface {
	Open() error
	Close() error
	Read(int, int) (string, error)
	Write(int, string) error
	Append(string) error
	ReadNextLine(int, int) (string, error)
}
