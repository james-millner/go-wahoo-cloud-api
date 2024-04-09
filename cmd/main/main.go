package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goji "goji.io"
	"io"
	"log"
	"net/http"
	"net/url"
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
}

// Config struct for holding environment variables.
type Config struct {
	HTTPPort int  `default:"8080"`
	Debug    bool `default:"false"`
}

// main function
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

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
	router.HandleFunc(pat.Post("/callback"), Callback())
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

	log.Println("Authorize called")
	//Redirect to Wahoo API
	redirectUrl, err := getWahooAuthorizeUrl(wahooClientId, wahooRedirectUri)
	log.Println(redirectUrl.String())

	return func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, redirectUrl.String(), 301)
	}
}

func getWahooAuthorizeUrl(wahooClientId, wahooRedirectUri string) (*url.URL, error) {
	return url.Parse("https://api.wahooligan.com/oauth/authorize?" +
		"client_id=" + wahooClientId +
		"&redirect_uri=" + wahooRedirectUri +
		"&scope=user_read%20workouts_read%20offline_data" +
		"&response_type=code")
}

func HomeRoot(wahooClientId, wahooClientSecret, wahooRedirectUri string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		log.Println("HomeRoot called")

		code := r.URL.Query().Get("code")

		if code == "" {
			log.Printf("No code found in the URL")
			http.Redirect(w, r, "https://api.wahooligan.com/oauth/authorize?"+
				"client_id="+wahooClientId+
				"&redirect_uri="+wahooRedirectUri+
				"&scope=user_read%20workouts_read%20offline_data"+
				"&response_type=code", 301)
		}

		oauthUrl, err := getWahooOAuthExchangeURL(wahooClientId, wahooClientSecret, code, wahooRedirectUri)

		log.Printf("OAuth URL: %s", oauthUrl.String())
		log.Printf("code: %s", code)

		if err != nil {
			log.Printf("Error getting the OAuth exchange URL: %v", err)
			panic(err)
		}

		resp, err := http.Post(oauthUrl.String(), "application/json", nil) // "application/x-www-form-urlencoded
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

		log.Println(string(body))

		if resp.StatusCode == http.StatusOK {
			log.Printf("Authorization code received: %s", string(body))
			// Respond with the authorization code if needed.
			fprintf, err := fmt.Fprintf(w, string(body))
			if err != nil {
				log.Fatal(fprintf)
				return
			}

			resp := &Response{
				Message: "OK",
			}

			enc.Encode(resp)

		} else {
			log.Printf("Response status code: %d", resp.StatusCode)
			log.Printf("Error response: %s", string(body))
			// Respond with an error message if needed.
			fprintf, err := fmt.Fprintf(w, string(body))
			if err != nil {
				log.Fatal(fprintf)
				return
			}

			resp := &Response{
				Message: "Error Response Whoops",
			}

			enc.Encode(resp)
		}
	}
}

func getWahooOAuthExchangeURL(wahooClientId, wahooClientSecret, code, wahooRedirectUri string) (*url.URL, error) {
	return url.Parse("https://api.wahooligan.com/oauth/token?" +
		"client_id=" + wahooClientId +
		"&client_secret=" + wahooClientSecret +
		"&code=" + code +
		"&grant_type=authorization_code" +
		"&redirect_uri=" + wahooRedirectUri)
}

func Callback() func(w http.ResponseWriter, r *http.Request) {

	log.Println("Callback called")

	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}

		log.Println("Request Body: ", string(body))

		resp := &Response{
			Message: string(body),
		}

		enc.Encode(resp)
	}
}
