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
	Address       string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_DURATION" envDefault:"10s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore from file")
	flag.StringVar(&cfg.Address, "a", cfg.Address, "host:port to listen")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "where to store data")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store interval")
	flag.Parse()
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
