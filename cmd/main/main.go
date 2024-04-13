package main

import (
	"context"
	"errors"
	goji "goji.io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/health"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/oauth"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/webhook"

	"goji.io/pat"
)

// Response Struct
type Response struct {
	Message string `json:"message"`
}

// Config struct for holding environment variables.
type Config struct {
	HTTPPort int  `default:"8080"`
	Debug    bool `default:"false"`
}

// main function
func main() {

	wahooClientId := os.Getenv("WAHOO_CLIENT_ID")
	wahooClientSecret := os.Getenv("WAHOO_CLIENT_SECRET")
	wahooRedirectURI := os.Getenv("REDIRECT_URI")
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(httpPort),
		Handler: handlersMethod(wahooClientId, wahooClientSecret, wahooRedirectURI),
	}

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

func handlersMethod(wahooClientId, wahooClientSecret, wahooRedirectUri string) *goji.Mux {
	router := goji.NewMux()

	router.HandleFunc(pat.Get("/healthz"), health.HealthHandler())
	router.HandleFunc(pat.Get("/authorize"), oauth.Authorize(wahooClientId, wahooRedirectUri))
	router.HandleFunc(pat.Get("/"), oauth.AuthCallback(wahooClientId, wahooClientSecret, wahooRedirectUri))
	router.HandleFunc(pat.Post("/callback"), webhook.Callback())
	return router
}
