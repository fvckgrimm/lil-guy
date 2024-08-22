package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func logError(err error) {
	logFile := filepath.Join(os.TempDir(), "lil-guy-error.log")
	f, openErr := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		fmt.Printf("Failed to open log file: %v\n", openErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Error: %v\n", err)

	fmt.Printf("Error logged to: %s\n", logFile)
}

func logOutput(output string) {
	logFile := filepath.Join(os.TempDir(), "lil-guy-output.log")
	f, openErr := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		fmt.Printf("Failed to open log file: %v\n", openErr)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(output)
}
