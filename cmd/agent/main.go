package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whiterale/go-prac/internal/metrics"
)

func main() {
	pollInterval := 1 * time.Second
	reportInterval := 5 * time.Second

	// start the tickers
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	// do not forget to stop them
	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
		log.Println("Bye")
	}()

	// take care of signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	metrics := metrics.Init()
	// main loop
	for {
		select {
		case <-pollTicker.C:
			log.Println("Polling")
			metrics.Poll()
		case <-reportTicker.C:
			log.Println("Reporting")
			metrics.Report("abc")
		case <-signals:
			log.Println("yal8r...")
			return
		}
	}
}
