package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cast"
	"github.com/ssyrota/frog-db/src/core/db"
	"github.com/ssyrota/frog-db/src/web"
	"github.com/tj/go/env"
)

func main() {
	if err := load(); err != nil {
		log.Fatal(err)
	}
}
func load() error {
	port := cast.ToInt64(env.GetDefault("PORT", "8080"))
	dumpPath := env.GetDefault("DUMP_PATH", "dump.json")
	dumpInterval, err := time.ParseDuration(env.GetDefault("DUMP_IVL", "1m"))
	if err != nil {
		return fmt.Errorf("parse dump ivl: %w", err)
	}
	db, err := db.New(dumpPath, dumpInterval)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	if err := web.New(db, uint16(port)).Run(); err != nil {
		return fmt.Errorf("run rest: %w", err)
	}
	return nil
}
