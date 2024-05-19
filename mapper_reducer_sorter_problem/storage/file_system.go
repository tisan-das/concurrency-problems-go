package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type FileSystem struct {
	fileName          string
	openedFile        *os.File
	openedFileScanner *bufio.Scanner
}

func NewFileSystem(fileName string) File {
	return &FileSystem{
		fileName:          fileName,
		openedFile:        nil,
		openedFileScanner: nil,
	}
}

func (file *FileSystem) Open() error {
	// no-op
	return nil
}

func (file *FileSystem) Close() error {
	// no-op
	file.openedFile.Close()
	return nil
}

func (file *FileSystem) Read(int, int) (string, error) {
	openedFile, err := os.OpenFile(file.fileName, os.O_RDONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("Error occurred while opening file %s: %s", file.fileName, err)
	}
	defer openedFile.Close()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(openedFile)
	if err != nil {
		return "", fmt.Errorf("Error occurred while writing the content to file %s : %s", file.fileName, err)
	}
	return buffer.String(), nil

}

func (file *FileSystem) ReadNextLine(int, int) (string, error) {
	if file.openedFile == nil {
		openedFile, err := os.OpenFile(file.fileName, os.O_RDONLY, 0600)
		if err != nil {
			return "", fmt.Errorf("Error occurred while opening file %s: %s", file.fileName, err)
		}
		file.openedFile = openedFile
		scanner := bufio.NewScanner(openedFile)
		file.openedFileScanner = scanner
	}
	if file.openedFileScanner.Scan() {
		return file.openedFileScanner.Text(), nil
	} else {
		return "", nil
	}
}

func (file *FileSystem) Write(offset int, content string) error {
	openedFile, err := os.OpenFile(file.fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Error occurred while opening file %s: %s", file.fileName, err)
	}
	defer openedFile.Close()
	buffer := []byte(content)
	_, err = openedFile.WriteAt(buffer, int64(offset))
	if err != nil {
		return fmt.Errorf("Error occurred while writing to file %s: %s", file.fileName, err)
	}
	return nil
}

func (file *FileSystem) Append(content string) error {

	if file.openedFile == nil {
		openedFile, err := os.OpenFile(file.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return fmt.Errorf("Error occurred while opening file %s: %s", file.fileName, err)
		}
		file.openedFile = openedFile
	}
	_, err := file.openedFile.WriteString(content)
	if err != nil {
		return fmt.Errorf("Error occurred while writing the content to file %s : %s", file.fileName, err)
	}
	return nil

	// openedFile, err := os.OpenFile(file.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// if err != nil {
	// 	return fmt.Errorf("Error occurred while opening file %s: %s", file.fileName, err)
	// }
	// defer openedFile.Close()
	// _, err = openedFile.WriteString(content)
	// if err != nil {
	// 	return fmt.Errorf("Error occurred while writing the content to file %s : %s", file.fileName, err)
	// }
	// return nil
}
