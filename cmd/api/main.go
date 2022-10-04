package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/nicholasasimov/findhotel-assignment/config"
	"github.com/nicholasasimov/findhotel-assignment/server"
	"github.com/nicholasasimov/findhotel-assignment/store"
)

func main() {
	ctx := context.Background()

	flag.Usage = config.Usage
	flag.Parse()

	if err := config.Init(); err != nil {
		fmt.Printf("can't init config: %s\n\n", err)
		config.Usage()
		os.Exit(1)
	}

	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	if config.Config.PrettyLog {
		log = log.Output(zerolog.NewConsoleWriter())
	}

	db, err := pgx.Connect(ctx, config.DSN())
	if err != nil {
		log.Error().Err(err).Msg("can't connect to database")
		return
	}
	defer db.Close(ctx)

	// TODO migrate func to create table
	s := server.New(log, db, store.New(db))

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer, middleware.RealIP, middleware.RequestID, middleware.Compress(5))

	r.Get("/health", Health)
	r.Mount("/debug", middleware.Profiler())
	r.Mount("/api/v1", s.Routes())

	srv := &http.Server{
		Addr:              net.JoinHostPort(config.Config.Server.Addr, config.Config.Server.Port),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		Handler:           r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info().Msg("received signal, stopping server..")
		_ = srv.Shutdown(ctx)
	}()

	log.Info().Str("addr", srv.Addr).Msg("start server")
	if err := srv.ListenAndServe(); err != nil {
		log.Info().Err(err).Msg("server stopped")
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
