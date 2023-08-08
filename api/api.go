package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"github.com/twitter-remake/user/api/middlewares"
	"github.com/twitter-remake/user/backend"
	"github.com/twitter-remake/user/config"
	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

type JSON map[string]any

type Server struct {
	app *http.Server
}

// New creates a new user twirp server
func New(backend *backend.Dependency) *Server {
	handler := newHandler(backend)
	listingServiceServer := userpb.NewListingServer(handler)
	profileServiceServer := userpb.NewProfileServer(handler)

	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept"},
	}))
	mux.Use(middleware.RequestID)
	mux.Use(middlewares.Helmet)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ReplyJSON(w, http.StatusOK, JSON{
			"status": "OK",
		})
	})
	mux.Mount(listingServiceServer.PathPrefix(), listingServiceServer)
	mux.Mount(profileServiceServer.PathPrefix(), profileServiceServer)

	httpsrv := &http.Server{
		Addr:         net.JoinHostPort(config.Host(), config.Port()),
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	return &Server{
		app: httpsrv,
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

func ReplyJSON(w http.ResponseWriter, status int, body JSON) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type Handler interface {
	userpb.Listing
	userpb.Profile
	// userpb.Relationship
}

type handler struct {
	backend *backend.Dependency
}

func newHandler(backend *backend.Dependency) Handler {
	return &handler{backend: backend}
}
