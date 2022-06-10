package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/whiterale/go-prac/internal/server"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%+v\n", cfg)
	server.Listen(cfg.Address)
}
