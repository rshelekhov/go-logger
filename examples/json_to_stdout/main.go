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

	// Imitation of logging from multiple goroutines
	go func() {
		for i := 0; i < 5; i++ {
			log.Debug(fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(logger.DEBUG), i))
			time.Sleep(200 * time.Millisecond)

			log.Info(fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(logger.INFO), i))
			time.Sleep(200 * time.Millisecond)

			log.Warning(fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(logger.WARNING), i))
			time.Sleep(200 * time.Millisecond)

			log.Error(fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(logger.ERROR), i))
			time.Sleep(200 * time.Millisecond)

			log.Fatal(fmt.Sprintf("Goroutine 1 – %s message %d", logger.GetLevelString(logger.FATAL), i))
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			log.Debug(fmt.Sprintf("Goroutine 2 – %s message %d", logger.GetLevelString(logger.DEBUG), i))
			time.Sleep(300 * time.Millisecond)

			log.Info(fmt.Sprintf("Goroutine 2 – %s message %d", logger.GetLevelString(logger.INFO), i))
			time.Sleep(300 * time.Millisecond)

			log.Warning(fmt.Sprintf("Goroutine 2 – %s message %d", logger.GetLevelString(logger.WARNING), i))
			time.Sleep(300 * time.Millisecond)

			log.Error(fmt.Sprintf("Goroutine 2 – %s message %d", logger.GetLevelString(logger.ERROR), i))
			time.Sleep(300 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
}
