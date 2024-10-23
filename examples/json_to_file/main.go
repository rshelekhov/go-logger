package main

import (
	"fmt"
	"github.com/rshelekhov/go-logger"
	"os"
	"time"
)

func main() {
	// Create log file
	file, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	// Create new logger instance
	log := logger.New(logger.DEBUG, file, true)
	defer log.Close()

	// Make slice of log levels for using in example
	levels := []int{logger.DEBUG, logger.INFO, logger.WARNING, logger.ERROR}

	// Imitation of logging from multiple goroutines
	go func() {
		for i := 0; i < 5; i++ {
			for level := range levels {
				log.Log(levels[level], fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(levels[level]), i))
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			for level := range levels {
				log.Log(levels[level], fmt.Sprintf("Goroutine 2 – %s message %d", logger.GetLevelString(levels[level]), i))
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	time.Sleep(5 * time.Second)
}
