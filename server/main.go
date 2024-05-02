package main

import (
	"context"
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
	context := context.Background()

	redisAddr := ""
	switch configuration.Environment {
	case "development":
		redisAddr = configuration.DevRedisAddress
	case "production":
		redisAddr = configuration.ProdRedisAddress
	case "testing":
		redisAddr = configuration.TestRedisAddress
	default:
		log.Fatal().Msg("unknown environment")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "changeme",
		DB:       0, // default db
	})

	if err := client.Ping(context).Err(); err != nil {
		log.Fatal().Err(err).Msg("cannot connect to Redis server")
	}

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
