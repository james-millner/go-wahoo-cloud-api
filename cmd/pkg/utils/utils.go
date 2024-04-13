package utils

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

func CheckIfAuthCodeDoesntExist(w http.ResponseWriter, r *http.Request, code string, wahooClientId string, wahooRedirectUri string) bool {
	if code == "" {
		log.Printf("No code found in the URL")
		http.Redirect(w, r, os.Getenv("WAHOO_AUTH_BASE_URL")+"?"+
			"client_id="+wahooClientId+
			"&redirect_uri="+wahooRedirectUri+
			"&scope=user_read%20workouts_read%20offline_data"+
			"&response_type=code", 301)
		return true
	}
	return false
}

func GetWahooOAuthExchangeURL(wahooClientId, wahooClientSecret, code, wahooRedirectUri string) (*url.URL, error) {
	tokenUrl := os.Getenv("WAHOO_TOKEN_BASE_URL")
	return url.Parse(tokenUrl + "?" +
		"client_id=" + wahooClientId +
		"&client_secret=" + wahooClientSecret +
		"&code=" + code +
		"&grant_type=authorization_code" +
		"&redirect_uri=" + wahooRedirectUri)
}

func GetWahooAuthorizeUrl(wahooClientId, wahooRedirectUri string) (*url.URL, error) {

	authUrl := os.Getenv("WAHOO_AUTH_BASE_URL")
	return url.Parse(authUrl + "?" +
		"client_id=" + wahooClientId +
		"&redirect_uri=" + wahooRedirectUri +
		"&scope=user_read%20workouts_read%20offline_data" +
		"&response_type=code")
}
