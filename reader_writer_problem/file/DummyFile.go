package file

import (
	"fmt"
	"sync"
)

type DummyFile struct {
	Location         string
	readers, writers int
}

var mutex sync.Mutex
var files map[string]bool = make(map[string]bool)

func NewDummyFile(location string) (File, error) {
	mutex.Lock()
	defer mutex.Unlock()
	value, exists := files[location]
	if exists && value {
		return &DummyFile{}, fmt.Errorf("File %s is already opened", location)
	}
	return &DummyFile{
		Location: location,
		readers:  0,
		writers:  0,
	}, nil
}

func (file *DummyFile) Open() error {
	mutex.Lock()
	defer mutex.Unlock()
	value, exists := files[file.Location]
	if exists && value {
		return fmt.Errorf("File %s is already opened", file.Location)
	}
	files[file.Location] = true
	return nil
}

func (file *DummyFile) Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	value, exists := files[file.Location]
	if !exists || !value {
		return fmt.Errorf("Please check whether the file %d was opened!", file.Location)
	}
	files[file.Location] = false
	return nil
}

func (file *DummyFile) Read(offset, bytes int) (string, error) {
	value, exists := files[file.Location]
	if !exists || !value {
		return "", fmt.Errorf("Please check whether the file %d was opened!", file.Location)
	}
	return "sample-text", nil
}

func (file *DummyFile) Write(offset int, content string) error {
	value, exists := files[file.Location]
	if !exists || !value {
		return fmt.Errorf("Please check whether the file %d was opened!", file.Location)
	}
	return nil
}

func (file *DummyFile) GetMetadata() string {
	// value, exists := files[file.Location]
	// if !exists || !value {
	// 	return "", fmt.Errorf("Please check whether the file %d was opened!", file.Location)
	// }
	return file.Location
}
