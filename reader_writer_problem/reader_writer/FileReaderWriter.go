package reader_writer

import (
	"concurrency_problems/reader_writer_problem/file"
	"log"
	"sync"
	"time"
)

type FileReaderWriter struct {
	file             file.File
	readers, writers int
	writerLock       sync.Mutex
	readerLock       sync.Mutex
}

func NewFileReaderWriter(file file.File) ReaderWriter {
	return &FileReaderWriter{
		file:    file,
		readers: 0,
		writers: 0,
	}
}

func (readerWriter *FileReaderWriter) Open() error {
	return readerWriter.file.Open()
}

func (readerWriter *FileReaderWriter) Close() error {
	return readerWriter.file.Close()
}

func (readerWriter *FileReaderWriter) Read(offset, bytes int) (string, error) {
	// Set Writer lock for the first reading activity
	readerWriter.readerLock.Lock()
	if readerWriter.readers == 0 {
		readerWriter.writerLock.Lock()
	}
	readerWriter.readers += 1
	readerWriter.readerLock.Unlock()

	// Read activity
	log.Printf("Current ReaderWriter perf while reading: (Readers:%d, Writers:%d)\n", readerWriter.readers, readerWriter.writers)
	time.Sleep(5 * time.Second)
	content, err := readerWriter.file.Read(offset, bytes)

	// Unset Writer lock for the last reading activity
	readerWriter.readerLock.Lock()
	if readerWriter.readers == 1 {
		readerWriter.writerLock.Unlock()
	}
	readerWriter.readers -= 1
	readerWriter.readerLock.Unlock()

	return content, err
}

func (readerWriter *FileReaderWriter) Write(offset int, content string) error {
	// Obtain writer lock to allow only one writer
	readerWriter.writerLock.Lock()
	defer readerWriter.writerLock.Unlock()

	// Write Activity
	readerWriter.writers += 1
	log.Printf("Current ReaderWriter perf while writing: (Readers:%d, Writers:%d)\n", readerWriter.readers, readerWriter.writers)
	time.Sleep(10 * time.Second)
	err := readerWriter.file.Write(offset, content)
	readerWriter.writers -= 1
	return err
}

func (readerWriter *FileReaderWriter) GetFileDetails() string {
	return readerWriter.file.GetMetadata()
}
