package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
)

type WahooTokenResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	TokenType    string `json:"token_type" validate:"required"`
	ExpiresIn    int    `json:"expires_in" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
	Scope        string `json:"scope" validate:"required"`
	CreatedAt    int    `json:"created_at" validate:"required"`
}

func Authorize() func(w http.ResponseWriter, r *http.Request) {

	log.Println("Authorize called")
	wahooClientId := os.Getenv("WAHOO_CLIENT_ID")
	wahooRedirectUri := os.Getenv("REDIRECT_URI")

	//Redirect to Wahoo API
	redirectUrl, err := utils.GetWahooAuthorizeUrl(wahooClientId, wahooRedirectUri)
	log.Println(redirectUrl.String())

	return func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, redirectUrl.String(), 301)
	}
}

func AuthCallback() func(w http.ResponseWriter, r *http.Request) {

	wahooClientId := os.Getenv("WAHOO_CLIENT_ID")
	wahooClientSecret := os.Getenv("WAHOO_CLIENT_SECRET")
	wahooRedirectUri := os.Getenv("REDIRECT_URI")

	return func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")

		if utils.CheckIfAuthCodeDoesntExist(w, r, code, wahooClientId, wahooRedirectUri) {
			return
		}

		oauthUrl, err := utils.GetWahooOAuthExchangeURL(wahooClientId, wahooClientSecret, code, wahooRedirectUri)
		if err != nil {
			log.Printf("Error getting the OAuth exchange URL: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		oauthResponse, err := http.Post(oauthUrl.String(), "application/json", nil) // "application/x-www-form-urlencoded
		if err != nil {
			log.Printf("Error making the POST request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if oauthResponse.StatusCode != http.StatusOK {
			log.Printf("Error response status code: %d", oauthResponse.StatusCode)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer oauthResponse.Body.Close()

		body, err := io.ReadAll(oauthResponse.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		if oauthResponse.StatusCode == http.StatusOK {
			fmt.Println("OAuth exchange successful. Response body:", string(body))
			var tokenResponse WahooTokenResponse
			jErr := json.Unmarshal(body, &tokenResponse)
			if jErr != nil {
				fmt.Println("Error unmarshalling JSON:", jErr)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			tokenValidator := validator.New(validator.WithRequiredStructEnabled())
			err = tokenValidator.Struct(tokenResponse)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", jErr)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			enc.Encode(tokenResponse)
			return
		}

		log.Printf("Response status code: %d", oauthResponse.StatusCode)
		log.Printf("Error response: %s", string(body))
		// Respond with an error message if needed.
		fprintf, err := fmt.Fprintf(w, string(body))
		if err != nil {
			log.Fatal(fprintf)
			return
		}

		enc.Encode(body)
		return
	}
}
