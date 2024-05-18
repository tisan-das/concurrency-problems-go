package reader_writer

import (
	"concurrency_problems/reader_writer_problem/file"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type FileReaderWriterSemaphore struct {
	file                   file.File
	readers, writers       int
	writerLock, readerLock sync.Mutex
	readerSemaphore        *semaphore.Weighted
}

func NewFileReaderWriterSemaphore(file file.File, maxParallelReaders int) ReaderWriter {
	fmt.Println("---> ", maxParallelReaders)
	return &FileReaderWriterSemaphore{
		file:            file,
		readers:         0,
		writers:         0,
		writerLock:      sync.Mutex{},
		readerLock:      sync.Mutex{},
		readerSemaphore: semaphore.NewWeighted(int64(maxParallelReaders)),
	}
}

func (readerWriter *FileReaderWriterSemaphore) Open() error {
	return readerWriter.file.Open()
}

func (readerWriter *FileReaderWriterSemaphore) Close() error {
	return readerWriter.file.Close()
}

func (readerWriter *FileReaderWriterSemaphore) Read(offset, bytes int) (string, error) {
	// Restrict writer for the first reader
	readerWriter.readerSemaphore.Acquire(context.Background(), 1) // thread-safe?
	readerWriter.readerLock.Lock()
	if readerWriter.readers == 0 {
		readerWriter.writerLock.Lock()
	}
	readerWriter.readers += 1
	readerWriter.readerLock.Unlock()

	// Read operation
	log.Printf("Current ReaderWriter perf while reading: (Readers:%d, Writers:%d)\n",
		readerWriter.readers, readerWriter.writers)
	time.Sleep(5 * time.Second)
	value, err := readerWriter.file.Read(offset, bytes)

	// Release writer lock for the last reader
	readerWriter.readerLock.Lock()
	if readerWriter.readers == 1 {
		readerWriter.writerLock.Unlock()
	}
	readerWriter.readers -= 1
	readerWriter.readerLock.Unlock()
	readerWriter.readerSemaphore.Release(1)
	return value, err
}

func (readerWriter *FileReaderWriterSemaphore) Write(offset int, content string) error {
	// Obtain the writer lock
	readerWriter.writerLock.Lock()
	defer readerWriter.writerLock.Unlock()

	// Write operation
	readerWriter.writers += 1
	log.Printf("Current ReaderWriter perf while writing: (Readers:%d, Writers:%d)\n",
		readerWriter.readers, readerWriter.writers)
	time.Sleep(10 * time.Second)
	err := readerWriter.file.Write(offset, content)
	readerWriter.writers -= 1
	return err
}

func (readerWriter *FileReaderWriterSemaphore) GetFileDetails() string {
	return readerWriter.file.GetMetadata()
}
