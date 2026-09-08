package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/durgasatkar/olx-api/internal/config"
	"github.com/durgasatkar/olx-api/internal/db"
	"github.com/durgasatkar/olx-api/internal/handlers"
	"github.com/durgasatkar/olx-api/internal/middleware"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	mux := http.NewServeMux()
	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.GetListings)
	mux.HandleFunc("DELETE /listings/{id}", lh.DeleteListing)

	handler := middleware.RequestId(mux)
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("database connected")
	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
