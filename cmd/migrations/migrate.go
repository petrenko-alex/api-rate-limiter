package main

import (
	"flag"
	"log"
	"os"

	_ "api-rate-limiter/migrations" // migrations import
	_ "github.com/lib/pq"           // driver import
	"github.com/pressly/goose/v3"
)

const DSN = "postgresql://main:main@localhost:5432/rate_limiter?sslmode=disable" // todo: to config
const MigrationsDir = "./migrations/"

func main() {
	os.Exit(run())
}

func run() int {
	args := os.Args[1:]
	if len(args) < 1 {
		flag.Usage()
		return 0
	}
	command := args[0]

	db, err := goose.OpenDBWithDriver("postgres", DSN)
	if err != nil {
		log.Printf("goose: failed to open DB: %v\n\n", err)

		return 1
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.Run(command, db, MigrationsDir, arguments...); err != nil {
		log.Printf("goose %v: %v\n", command, err)

		return 1
	}

	return 0
}
