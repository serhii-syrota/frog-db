package main

import (
	"fmt"
	"log"
	"time"

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
	port := 8080
	dumpPath := env.Get("DUMP_PATH")
	dumpInterval, err := time.ParseDuration(env.Get("DUMP_IVL"))
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
