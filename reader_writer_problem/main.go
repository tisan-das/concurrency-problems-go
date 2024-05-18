package main

import (
	"concurrency_problems/reader_writer_problem/file"
	"concurrency_problems/reader_writer_problem/reader_writer"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello!")

	file1, err := file.NewDummyFile("sample.txt")
	if err != nil {
		fmt.Println("Error occurred while creating file: ", err)
		return
	}

	file2, err := file.NewDummyFile("sample.txt")
	if err != nil {
		fmt.Println("Error occurred while creating file: ", err)
		return
	}

	var readerWriterFile1, readerWriterFile2 reader_writer.ReaderWriter

	// Example: ReaderWriter with unlimited parallel readers
	readerWriterFile1 = reader_writer.NewFileReaderWriter(file1)
	readerWriterFile2 = reader_writer.NewFileReaderWriter(file2)

	// Example: ReaderWriter with sempahore to restrict no of simultaneous readers
	// readerWriterFile1 = reader_writer.NewFileReaderWriterSemaphore(file1, 5)
	// readerWriterFile2 = reader_writer.NewFileReaderWriterSemaphore(file2, 5)

	err = readerWriterFile1.Open()
	if err != nil {
		fmt.Println("Error ocurred while opening file: ", readerWriterFile1.GetFileDetails())
	}

	err = readerWriterFile2.Open()
	if err != nil {
		fmt.Println("Error ocurred while opening file: ", readerWriterFile2.GetFileDetails(),
			" --> this one is expected")
	}

	// reader_writer.NewFileReaderWriter(file)
	fmt.Println("---> ", readerWriterFile1, readerWriterFile2)

	// Use threading for simultaneous operation on file
	log.Print("Initiating Simultaneous Reader Writer Execution!")
	var wGroup sync.WaitGroup
	for i := 0; i <= 10; i++ {
		wGroup.Add(1)
		go readThread(readerWriterFile1, i, &wGroup)
	}
	for i := 0; i <= 5; i++ {
		wGroup.Add(1)
		go writeThread(readerWriterFile1, i, &wGroup)
	}
	wGroup.Wait()
	log.Println("Execution Completed!!")
}

func readThread(readerWriter reader_writer.ReaderWriter, index int, wGroup *sync.WaitGroup) {
	for i := 0; i <= 5; i++ {
		log.Print("Triggering a reader call from reader thread:", index, " itr:", i)
		time.Sleep(1 * time.Second)
		_, err := readerWriter.Read(1, 1)
		if err != nil {
			fmt.Printf("Error occurred while reading from thread %d: %s\n", index, err)
		}
	}
	wGroup.Done()
}

func writeThread(readerWriter reader_writer.ReaderWriter, index int, wGroup *sync.WaitGroup) {
	for i := 0; i <= 5; i++ {
		log.Print("Triggering a writer call from reader thread:", index, " itr:", i)
		time.Sleep(1 * time.Second)
		err := readerWriter.Write(1, "")
		if err != nil {
			fmt.Printf("Error occurred while writing from thread %d: %s\n", index, err)
		}
	}
	wGroup.Done()
}
