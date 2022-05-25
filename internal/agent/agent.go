package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whiterale/go-prac/internal/metric"
)

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration

	report     interface{}
	buffer     Updater
	collectors []Collector
}

type Collector interface {
	Collect() []AgentMetric
}

type Updater interface {
	Update(string, string, interface{}) error
	Dump() []*metric.Metric
	Flush()
}

func Init(poll time.Duration, report time.Duration, reporter interface{}, buffer Updater) *Agent {
	return &Agent{poll, report, reporter, buffer, nil}
}

type AgentMetric struct {
	ID    string
	MType string
	Value interface{}
}

func (agent *Agent) Poll() {
	for _, c := range agent.collectors {
		metrics := c.Collect()
		for _, metric := range metrics {
			agent.buffer.Update(metric.MType, metric.ID, metric.Value)
		}
	}
}

func (agent *Agent) Report() {}

func (agent *Agent) Start() error {

	pollTicker := time.NewTicker(agent.pollInterval)
	reportTicker := time.NewTicker(agent.reportInterval)

	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
		log.Println("Bye")
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-pollTicker.C:
			log.Println("poll...")
			for _, collector := range agent.collectors {
				metrics := collector.Collect()
				for _, metric := range metrics {
					agent.buffer.Update(metric.MType, metric.ID, metric.Value)
				}
			}
		case <-reportTicker.C:
			for _, m := range agent.buffer.Dump() {
				fmt.Printf("%s\n", m)
			}
			agent.buffer.Flush()

		case <-signals:
			log.Println("stop.")
			return nil
		}
	}
}
