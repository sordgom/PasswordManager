package main

import (
	"os"

	"github.com/sordgom/PasswordManager/server/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	runServer(config)
}

func runServer(config config.Config) {
	log.Info().Str("address", config.ServerAddress).Msg("starting server")
}
