package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	const (
		numProcesses = 50
		duration     = 50 * time.Second // Change this to your desired timeout
		logDir       = "logs"
	)

	var wg sync.WaitGroup

	// Create log directory if not exists
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	for i := 0; i < numProcesses; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			logFilePath := filepath.Join(logDir, fmt.Sprintf("process_%d.log", id))
			logFile, err := os.Create(logFilePath)
			if err != nil {
				fmt.Printf("Failed to create log file for process %d: %v\n", id, err)
				return
			}
			defer logFile.Close()

			fmt.Printf("Starting subprocess %d, logging to %s\n", id, logFilePath)

			// Create a context that times out after 'duration'
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			// Build command with context and args
			cmd := exec.CommandContext(ctx, "go", "run", "../.", "--doc-id=38a14415-8cb3-4d70-ab9f-e4d1a6463e7f", "--auto=true")

			cmd.Stdout = logFile
			cmd.Stderr = logFile

			// Start and wait for command
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(logFile, "Failed to start process %d: %v\n", id, err)
				return
			}

			if err := cmd.Wait(); err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					fmt.Fprintf(logFile, "Process %d killed after timeout of %v\n", id, duration)
				} else {
					fmt.Fprintf(logFile, "Process %d exited with error: %v\n", id, err)
				}
			} else {
				fmt.Fprintf(logFile, "Process %d completed successfully.\n", id)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All subprocesses completed.")
}
