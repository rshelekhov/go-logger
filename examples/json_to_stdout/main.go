package main

import (
	"fmt"
	"github.com/rshelekhov/go-logger"
	"os"
	"time"
)

func main() {
	// Create new logger instance
	log := logger.New(logger.DEBUG, os.Stdout, true)
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
