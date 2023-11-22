package main

import (
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq" // driver import
	"github.com/petrenko-alex/api-rate-limiter/internal/config"
	_ "github.com/petrenko-alex/api-rate-limiter/migrations" // migrations import
	"github.com/pressly/goose/v3"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	args := os.Args[1:]
	if len(args) < 2 {
		flag.Usage()
		return 0
	}
	command := args[1] //nolint:ifshort

	configFile, fileErr := os.Open(configFilePath)
	if fileErr != nil {
		log.Println("Error opening config file: ", fileErr)

		return 1
	}

	cfg, configErr := config.ForMigrator(context.Background(), configFile)
	if configErr != nil {
		log.Println("Error parsing config file: ", configErr)

		return 1
	}

	db, err := goose.OpenDBWithDriver("postgres", cfg.DB.DSN)
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
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, cfg.DB.MigrationsDir, arguments...); err != nil {
		log.Printf("goose %v: %v\n", command, err)

		return 1
	}

	return 0
}
