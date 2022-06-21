package main

import (
	"flag"
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
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "host:port to listen")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "store interval")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "store interval")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
		return
	}
	log.Printf("%+v\n", cfg)

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
