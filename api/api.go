package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/twitter-remake/user/backend"
	"github.com/twitter-remake/user/config"
	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

type Server struct {
	app *http.Server
}

// New creates a new user twirp server
func New(backend *backend.Dependency) *Server {
	handler := newHandler(backend)
	listingServiceServer := userpb.NewListingServer(handler)
	profileServiceServer := userpb.NewProfileServer(handler)

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	mux.Mount(listingServiceServer.PathPrefix(), listingServiceServer)
	mux.Mount(profileServiceServer.PathPrefix(), profileServiceServer)

	return &Server{
		app: &http.Server{
			Addr:    net.JoinHostPort(config.Host(), config.Port()),
			Handler: mux,
		},
	}
}

// Start starts the API server
func (s *Server) Start(host, port string) <-chan os.Signal {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		address := net.JoinHostPort(host, port)
		log.Info().Msgf("Listening on %s", address)
		err := s.app.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	return exitSignal
}

// Shutdown gracefully shuts down the API server
func (s *Server) Shutdown(ctx context.Context, signal os.Signal) {
	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	shutdownChan := make(chan error, 1)

	go func() {
		log.Warn().Any("signal", signal.String()).Msg("received signal, shutting down...")
		shutdownChan <- s.app.Shutdown(ctx)
	}()

	select {
	case <-timeout.Done():
		log.Warn().Msg("shutdown timed out, forcing exit")
		os.Exit(1)
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal().Err(err).Msg("there was an error shutting down")
		} else {
			log.Info().Msg("shutdown complete")
		}
	}
}

type Handler interface {
	userpb.Listing
	userpb.Profile
}

type handler struct {
	backend *backend.Dependency
}

func newHandler(backend *backend.Dependency) Handler {
	return &handler{backend: backend}
}
