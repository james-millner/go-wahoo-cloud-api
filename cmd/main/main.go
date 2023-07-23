package main

import (
	"context"
	"encoding/json"
	"fmt"
	goji "goji.io"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"goji.io/pat"
)

// Response Struct
type Response struct {
	Message string `json:"message"`
	Request string `json:"request"`
}

// Config struct for holding environment variables.
type Config struct {
	HTTPPort int  `default:"8092"`
	Debug    bool `default:"false"`
}

func main() {

	var env Config
	err := envconfig.Process("api", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Listening on: ", env.HTTPPort)

	wahooClientId := os.Getenv("WAHOO_CLIENT_ID")
	wahooClientSecret := os.Getenv("WAHOO_CLIENT_SECRET")
	wahooRedirectURI := os.Getenv("REDIRECT_URI")

	log.Println("Wahoo Client ID: ", wahooClientId)
	log.Println("Wahoo Client Secret: ", wahooClientSecret)
	log.Println("Wahoo Redirect URI: ", wahooRedirectURI)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(env.HTTPPort),
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

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Printf("HTTP Server shutdown!")
	}
}

func handlersMethod(wahooClientId, wahooClientSecret, wahooRedirectUri string) *goji.Mux {
	router := goji.NewMux()

	router.HandleFunc(pat.Get("/healthz"), Health())
	router.HandleFunc(pat.Get("/authorize"), Authorize(wahooClientId, wahooRedirectUri))
	router.HandleFunc(pat.Get("/"), HomeRoot(wahooClientId, wahooClientSecret, wahooRedirectUri))
	router.HandleFunc(pat.Get("/callback"), Callback())
	return router
}

// Health endpoint
func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		resp := &Response{
			Message: "OK",
		}

		enc.Encode(resp)
	}
}

func Authorize(wahooClientId, wahooRedirectUri string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Redirect to Wahoo API
		uri := fmt.Sprintf("https://api.wahooligan.com/oauth/authorize?"+
			"client_id=%s"+
			"&redirect_uri=%s"+
			"&scope=user_read workouts_read offline_data"+
			"&response_type=code",
			wahooClientId, wahooRedirectUri)

		http.Redirect(w, r, uri, 301)
	}
}

func HomeRoot(wahooClientId, wahooClientSecret, wahooRedirectUri string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		code := r.URL.Query().Get("code")

		accessTokenUrl := fmt.Sprintf("https://api.wahooligan.com/oauth/token?"+
			"client_id=%s"+
			"&client_secret=%s"+
			"&code=%s"+
			"&grant_type=authorization_code"+
			"&redirect_uri=%s",
			wahooClientId, wahooClientSecret, code, wahooRedirectUri)

		resp, err := http.Post(accessTokenUrl, "application/x-www-form-urlencoded", nil)
		if err != nil {
			log.Printf("Error making the POST request: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		if resp.StatusCode == http.StatusOK {
			log.Printf("Authorization code received: %s", string(body))
			// Respond with the authorization code if needed.
			fprintf, err := fmt.Fprintf(w, string(body))
			if err != nil {
				log.Fatal(fprintf)
				return
			}
		} else {
			log.Printf("Error response: %s", string(body))
			// Respond with an error message if needed.
			fprintf, err := fmt.Fprintf(w, string(body))
			if err != nil {
				log.Fatal(fprintf)
				return
			}
		}
	}
}

func Callback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		resp := &Response{
			Message: "OK",
			Request: string(body),
		}

		enc.Encode(resp)
	}
}
