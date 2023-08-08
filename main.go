package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/twitter-remake/user/api"
	"github.com/twitter-remake/user/backend"
	"github.com/twitter-remake/user/clients"
	"github.com/twitter-remake/user/config"
	"github.com/twitter-remake/user/repository"
)

func init() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Stack().Logger()
	if os.Getenv("ENVIRONMENT") == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting Twitter Auth Service")
	config.Init()
}

func main() {
	ctx := context.Background()

	// Initialize layers
	clients, err := clients.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize clients")
	}

	repository := repository.New(clients.PostgreSQL)
	backend := backend.New(clients, repository)
	api := api.New(backend)

	// Start server and wait for shutdown signals
	exitSignal := api.Start(config.Host(), config.Port())

	// If a shutdown signal is received (e.g. CTRL + C or kill) shutdown gracefully
	// signal stored in variable for logging purposes
	signal := <-exitSignal
	api.Shutdown(ctx, signal)
}
