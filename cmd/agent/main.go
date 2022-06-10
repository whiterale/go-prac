package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/whiterale/go-prac/internal/agent"
	"github.com/whiterale/go-prac/internal/agent/buffer"
	"github.com/whiterale/go-prac/internal/agent/collectors"
	"github.com/whiterale/go-prac/internal/agent/reporters"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%+v\n", cfg)

	agent := agent.Init(
		&reporters.JSON{Host: cfg.Address},
		buffer.Init(),
		[]agent.Collector{&collectors.Random{}, &collectors.PollCounter{}, &collectors.Runtime{}},
	)

	poll := time.NewTicker(cfg.PollInterval)
	report := time.NewTicker(cfg.ReportInterval)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		poll.Stop()
		report.Stop()
		log.Println("Bye")
	}()
	for {
		select {
		case <-poll.C:
			agent.Poll()
		case <-report.C:
			agent.Report()
		case <-signals:
			return
		}
	}
}
