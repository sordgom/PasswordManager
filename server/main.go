package main

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sordgom/PasswordManager/server/api"
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

func runServer(configuration config.Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     configuration.RedisAddress,
		Password: "changeme",
		DB:       0, // default db
	})

	vaultService := config.NewRedisVaultService(client)

	server, err := api.NewServer(configuration, vaultService)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	err = server.Start(configuration.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
