package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/whiterale/go-prac/internal/server"
	"github.com/whiterale/go-prac/internal/server/storage"
)

type config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func main() {
	cfg := config{}

	flag.BoolVar(&cfg.Restore, "r", true, "restore from file")
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "host:port to listen")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "where to store data")
	flag.DurationVar(&cfg.StoreInterval, "i", 10*time.Second, "store interval")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	log.Printf("%+v\n", cfg)

	// TODO: Naming looks ugly
	var st *storage.Storage
	if cfg.Restore {
		st, _ = storage.InitFromFile(cfg.StoreFile)
	} else {
		st = storage.Init()
	}

	if cfg.StoreInterval > 0 {
		stop := make(chan struct{})
		go st.StartSync(cfg.StoreFile, stop)
		defer func() {
			st.DumpToFile(cfg.StoreFile)
			stop <- struct{}{}
		}()
	}

	if cfg.StoreInterval == 0 {
		st.IsSync = true
	}

	// TODO: Naming looks ugly
	srv := server.Server{Storage: st}
	server.Listen(srv, cfg.Address)
}
