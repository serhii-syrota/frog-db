package db

import "time"

type config struct {
	Path         string        `env:"PATH" envDefault:"dump.json"`
	DumpInterval time.Duration `env:"DUMP_INTERVAL" envDefault:"1m"`
}
