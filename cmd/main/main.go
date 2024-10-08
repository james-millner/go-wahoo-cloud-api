package main

import (
	"context"
	"errors"
	goji "goji.io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/health"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/oauth"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/webhook"

	"goji.io/pat"
)

// main function
func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlersMethod(),
	}

	log.Printf("Starting server on port %v", port)

	go func() {
		// Graceful shutdown
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Unable to shut down server: %v", err)
		} else {
			log.Println("Server stopped")
		}
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%v", err)
	} else {
		log.Printf("HTTP Server shutdown!")
	}
}

func handlersMethod() *goji.Mux {
	router := goji.NewMux()

	router.HandleFunc(pat.Get("/healthz"), health.Health())
	router.HandleFunc(pat.Get("/"), oauth.AuthCallback())
	router.HandleFunc(pat.Post("/callback"), webhook.Callback())
	return router
}
