package main

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/ssyrota/frog-db/src/core/db"
	"github.com/ssyrota/frog-db/src/web"
)

type Config struct {
	Path         string        `env:"DUMP_PATH" envDefault:".dump.json"`
	DumpInterval time.Duration `env:"DUMP_INTERVAL" envDefault:"1m"`
	Port         int           `env:"REST_PORT" envDefault:"8080"`
}

func main() {
	if err := load(); err != nil {
		log.Fatal(err)
	}
}

func load() error {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("config: %w", err)
	}
	db, err := db.New(cfg.Path, cfg.DumpInterval)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	if err := web.New(db, uint16(cfg.Port)).Run(); err != nil {
		return fmt.Errorf("run rest: %w", err)
	}
	return nil
}
