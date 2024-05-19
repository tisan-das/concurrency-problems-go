package main

import (
	"concurrency_problems/mapper_reducer_sorter_problem/data_generator"
	"concurrency_problems/mapper_reducer_sorter_problem/storage"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	MAX_CONCURRENT_WRITER                              = 1000
	MAX_RANDOM_NUMBER_GENERATION_PER_INTERMEDIATE_FILE = 50000
	RANDOM_NUM_GENERATORS                              = 1000

	OPERATIONAL_DIR                 = "operations"
	UNSORTED_INITIAL_FILE_PREFIX    = "random_data_generated_"
	SORTED_INTERMEDIATE_FILE_PREFIX = "intermediate_sorted_data_"
	SORTED_FINAL_FILE_PREFIX        = "final_sorted_data"
	FILE_EXTENSION                  = ".txt"
)

var concurrentWriterSemaphore *semaphore.Weighted
var calculationMutex sync.Mutex
var totalInitialRandomNumberGenerated = 0
var totalEntriesInFinalSortedFile = 0

func main() {
	//fmt.Println("hello!")
	// testReadFile()
	concurrentWriterSemaphore = semaphore.NewWeighted(int64(MAX_CONCURRENT_WRITER))
	timeFormat := time.RFC3339
	startTime := time.Now()

	var SEED_VALUE int = -1
	generator := data_generator.NewRandomNumberGenerator(SEED_VALUE)
	var wGroup sync.WaitGroup

	dir, _ := os.Getwd()
	location := dir + "\\" + OPERATIONAL_DIR
	os.RemoveAll(location)
	os.Mkdir(location, 0777)

	// Part 1: Generate random numbers and store in file in unsorted way
	log.Print("------- Part 1: Generating random numbers -------")
	for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
		wGroup.Add(1)
		go fileWriterWithRandomNum(generator, i, &wGroup)
	}
	wGroup.Wait()
	log.Print("Initial random numbers generated: ", totalInitialRandomNumberGenerated)

	// Part 2: Mapper: Sort the content of individual files
	log.Print("------- Part 2: Mapper intermediate sorter -------")
	for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
		wGroup.Add(1)
		go intermediateSorterIndividualFile(i, &wGroup)
	}
	wGroup.Wait()

	// Part 3: Reducer: Capture and sort all the intermediate files in a streaming manner
	log.Print("------- Part 3: Reducer final sorter -------")
	finalSorter()
	log.Print("------- Completed -------")

	endTime := time.Now()
	fmt.Println("Execution started at: ", startTime.Format(timeFormat))
	fmt.Println("Execution ended at: ", endTime.Format(timeFormat))
	fmt.Println("Execution duration: ", endTime.Sub(startTime))
	fmt.Println("Initial random numbers generated: ", totalInitialRandomNumberGenerated)
	fmt.Println("Total entries at final sorted file: ", totalEntriesInFinalSortedFile)
}

/*
	Learnings:
		1. Always convert the integer to string using strconv.Itoa()
		2. What is range??
*/

func fileWriterWithRandomNum(generator *data_generator.RandomGenerator, index int, wGroup *sync.WaitGroup) {
	concurrentWriterSemaphore.Acquire(context.Background(), 1)
	defer concurrentWriterSemaphore.Release(1)
	defer wGroup.Done()
	log.Printf("Initiating random data generation for index: %d", index)
	count := generator.GenerateRandomInt() % MAX_RANDOM_NUMBER_GENERATION_PER_INTERMEDIATE_FILE
	log.Printf("Generating %d numbers for index: %d", count, index)
	calculationMutex.Lock()
	totalInitialRandomNumberGenerated += count
	calculationMutex.Unlock()
	var content strings.Builder
	for i := 0; i < count; i++ {
		content.WriteString(strconv.Itoa(generator.GenerateRandomInt()) + "\n")
	}

	// TODO: Handle directory structure in OS-agnostic way
	dir, _ := os.Getwd()
	fileName := dir + "\\" + OPERATIONAL_DIR + "\\" + UNSORTED_INITIAL_FILE_PREFIX + strconv.Itoa(index) + FILE_EXTENSION
	log.Printf("Storing the random data generated to file %s for index: %d", fileName, index)
	var file storage.File
	file = storage.NewFileSystem(fileName)
	err := file.Write(0, content.String())
	file.Close()
	if err != nil {
		log.Printf("Error occurred while writing content to file %s: %s", fileName, err)
	}
}

// TODO: use semaphore to restrict concurrent operations
func intermediateSorterIndividualFile(index int, wGroup *sync.WaitGroup) {
	concurrentWriterSemaphore.Acquire(context.Background(), 1)
	defer concurrentWriterSemaphore.Release(1)
	defer wGroup.Done()

	// Step 01: Read the content from the unsorted file
	dir, _ := os.Getwd()
	intermediateUnsortedFilename := dir + "\\" + OPERATIONAL_DIR + "\\" + UNSORTED_INITIAL_FILE_PREFIX + strconv.Itoa(index) + FILE_EXTENSION
	log.Printf("Reading the intermediate file %s for index: %d", intermediateUnsortedFilename, index)
	var file storage.File
	file = storage.NewFileSystem(intermediateUnsortedFilename)
	content, err := file.Read(-1, -1)
	file.Close()
	if err != nil {
		log.Printf("Error occurred while reading content from file %s: %s",
			intermediateUnsortedFilename, err)
	}

	// Step 02: Convert the content to int array and sort
	log.Printf("Converting and sorting the data read from file %s for index: %d",
		intermediateUnsortedFilename, index)
	contents := strings.Split(content, "\n")
	var array []int
	for _, value := range contents {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			array = append(array, intValue)
		}
	}
	sort.Ints(array)
	contents = make([]string, 0)
	content = ""

	var outputContent strings.Builder
	for _, value := range array {
		outputContent.WriteString(strconv.Itoa(value) + "\n")
	}

	// Step 03: Write the sorted content to intermediate sorted file
	outputFilename := dir + "\\" + OPERATIONAL_DIR + "\\" + SORTED_INTERMEDIATE_FILE_PREFIX + strconv.Itoa(index) + FILE_EXTENSION
	log.Printf("Storing the sorted data to file %s for index: %d", outputFilename, index)
	file = storage.NewFileSystem(outputFilename)
	err = file.Write(0, outputContent.String())
	file.Close()
	if err != nil {
		log.Printf("Error occurred while storing the sorted content to file %s for index: %d",
			outputFilename, index)
	}
}

func finalSorter() {
	log.Print("Initializing the cursor for all the intermediate sorted files")
	var files []storage.File
	dir, _ := os.Getwd()
	for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
		intermediateSortedFilename := dir + "\\" + OPERATIONAL_DIR + "\\" + SORTED_INTERMEDIATE_FILE_PREFIX + strconv.Itoa(i) + FILE_EXTENSION
		files = append(files, storage.NewFileSystem(intermediateSortedFilename))
	}

	// Set initial cursor for all the files
	var array [RANDOM_NUM_GENERATORS]string
	for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
		str, _ := files[i].ReadNextLine(-1, -1)
		array[i] = str
	}

	fmt.Println(array)

	finalOutputFileName := dir + "\\" + OPERATIONAL_DIR + "\\" + SORTED_FINAL_FILE_PREFIX + FILE_EXTENSION
	os.Remove(finalOutputFileName)
	var finalOutputFile storage.File
	finalOutputFile = storage.NewFileSystem(finalOutputFileName)

	log.Print("Initiating the streaming merge sort operation on all the opened file")
	var iteration = 0
	for {
		iteration++
		// Find the smallest element and increase corresponding cursor
		var indexWithMinValue int = -1
		var currentMinValue = 0
		for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
			if array[i] != "" && indexWithMinValue == -1 {
				indexWithMinValue = i
				value, _ := strconv.Atoi(array[i])
				currentMinValue = value
			}
			if indexWithMinValue != -1 && array[i] != "" {
				value, _ := strconv.Atoi(array[i])
				if value < currentMinValue {
					indexWithMinValue = i
					currentMinValue = value
				}
			}
		}
		// log.Printf("01 Iteration:%d Smallest Element:%d Index:%d", iteration, currentMinValue, indexWithMinValue)

		if indexWithMinValue == -1 {
			// Termination condition
			break
		} else {
			// Append it to the final sorted file
			finalOutputFile.Append(strconv.Itoa(currentMinValue) + "\n")

			// Increase the cursor of the minimal index
			str, _ := files[indexWithMinValue].ReadNextLine(-1, -1)
			array[indexWithMinValue] = str
		}

		// log.Printf("02 Iteration:%d Smallest Element:%d Index:%d", iteration, currentMinValue, indexWithMinValue)
	}
	totalEntriesInFinalSortedFile = iteration - 1

	// Close the files
	err := finalOutputFile.Close()
	if err != nil {
		log.Print("Error occurred while closing file %s: %s", finalOutputFileName, err)
	}
	for i := 0; i < RANDOM_NUM_GENERATORS; i++ {
		err := files[i].Close()
		if err != nil {
			log.Print("Error occurred while closing file for index %d: %s", i, err)
		}
	}
}

func testReadFile() {
	dir, _ := os.Getwd()
	intermediateSortedFilename := dir + "\\" + OPERATIONAL_DIR + "\\" + SORTED_INTERMEDIATE_FILE_PREFIX + "5" + FILE_EXTENSION
	file := storage.NewFileSystem(intermediateSortedFilename)
	fmt.Println(intermediateSortedFilename)
	for {
		var str string
		str, _ = file.ReadNextLine(-1, -1)
		if str == "" {
			break
		}
		fmt.Println(str)
		time.Sleep(time.Second)
	}
}
