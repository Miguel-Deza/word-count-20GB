package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func mainsds() {
	// Open the original file for reading
	originalFile, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	// Create the output file
	outputFile, err := os.Create("larger.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Get the size of the original file
	originalFileInfo, err := originalFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	originalSize := originalFileInfo.Size()
	fmt.Printf("original size: %v\n", originalSize)

	// Set the desired size of the output file (20 GB)
	desiredSize := int64(20) * 1024 * 1024 * 1024
	// 20 GB * 1024 MB/GB * 1024 KB/MB * 1024 bytes/KB
	// Un gigabyte equivale a 1024 megabytes (MB), y cada megabyte equivale a 1024 kilobytes (KB)

	// Calculate the number of times to replicate the content
	replicationFactor := desiredSize / originalSize

	sizeOfOutput := int64(0)

	// Replicate the content of the original file
	for i := int64(0); i < replicationFactor; i++ {
		// Set the file pointer to the beginning of the original file
		_, err := originalFile.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(outputFile, originalFile)
		if err != nil {
			log.Fatal(err)
		}
		sizeOfOutput += originalSize
	}

	// Truncate the output file to the desired size
	if sizeOfOutput > desiredSize {
		err = outputFile.Truncate(desiredSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("File generated successfully!")
}
