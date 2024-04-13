package utils

import "net/url"

func GetWahooOAuthExchangeURL(wahooClientId, wahooClientSecret, code, wahooRedirectUri string) (*url.URL, error) {
	return url.Parse("https://api.wahooligan.com/oauth/token?" +
		"client_id=" + wahooClientId +
		"&client_secret=" + wahooClientSecret +
		"&code=" + code +
		"&grant_type=authorization_code" +
		"&redirect_uri=" + wahooRedirectUri)
}

func GetWahooAuthorizeUrl(wahooClientId, wahooRedirectUri string) (*url.URL, error) {
	return url.Parse("https://api.wahooligan.com/oauth/authorize?" +
		"client_id=" + wahooClientId +
		"&redirect_uri=" + wahooRedirectUri +
		"&scope=user_read%20workouts_read%20offline_data" +
		"&response_type=code")
}
