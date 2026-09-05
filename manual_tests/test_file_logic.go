package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Simulate the loadFile function logic
func checkFile(filename string) error {
	// Check if file exists first
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Get file info for size check
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	// Check if file is too large (simulating 50MB limit)
	maxSizeBytes := int64(50) * 1024 * 1024
	fileSize := fileInfo.Size()
	
	// Convert file size to MB for display
	fileSizeMB := float64(fileSize) / (1024 * 1024)
	
	if fileSize > maxSizeBytes {
		return fmt.Errorf("file too large (%.1f MB), max %d MB", 
			fileSizeMB, 50)
	}

	// Try to read the file
	_, err = ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

func main() {
	// Test 1: Non-existent file
	fmt.Println("Test 1: Non-existent file")
	err := checkFile("nonexistent.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// This is the message that should appear in the dialog
		fmt.Printf("Dialog message would be: Cannot open file: nonexistent.txt\n%v\n", err)
	} else {
		fmt.Println("No error")
	}

	fmt.Println()

	// Test 2: Existing file (l.txt)
	fmt.Println("Test 2: Existing file (l.txt)")
	err = checkFile("l.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("No error - file loaded successfully")
	}
}