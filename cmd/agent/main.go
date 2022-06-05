package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whiterale/go-prac/internal/agent"
	"github.com/whiterale/go-prac/internal/agent/buffer"
	"github.com/whiterale/go-prac/internal/agent/collectors"
	"github.com/whiterale/go-prac/internal/agent/reporters"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	agent := agent.Init(
		&reporters.HTTPPlainText{Host: "http://localhost:8080"},
		buffer.Init(),
		[]agent.Collector{&collectors.Random{}, &collectors.PollCounter{}, &collectors.Runtime{}},
	)
	poll := time.NewTicker(1 * time.Second)
	report := time.NewTicker(3 * time.Second)

	defer func() {
		poll.Stop()
		report.Stop()
		log.Println("Bye")
	}()

	for {
		select {
		case <-poll.C:
			log.Println("Poll")
			agent.Poll()
		case <-report.C:
			log.Println("Report")
			agent.Report()
		case <-signals:
			return
		}
	}
}
