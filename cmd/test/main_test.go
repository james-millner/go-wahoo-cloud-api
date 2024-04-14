package test

import (
	"fmt"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/pkg/utils"
	"testing"
)

func TestGetWahooAuthorizeURL(t *testing.T) {

	t.Setenv("WAHOO_CLIENT_ID", "client123")
	t.Setenv("WAHOO_CLIENT_SECRET", "client_secret")
	t.Setenv("REDIRECT_URI", "https://example.com/callback")
	t.Setenv("WAHOO_AUTH_BASE_URL", "https://api.wahooligan.com/oauth/authorize")
	t.Setenv("WAHOO_TOKEN_BASE_URL", "https://api.wahooligan.com/oauth/token")

	testCases := []struct {
		name             string
		wahooClientId    string
		wahooRedirectUri string
		expectedResult   string
		expectedError    bool
	}{
		{
			name:             "Valid input",
			wahooClientId:    "client123",
			wahooRedirectUri: "https://example.com/callback",
			expectedResult:   "https://api.wahooligan.com/oauth/authorize?client_id=client123&redirect_uri=https://example.com/callback&scope=user_read%20workouts_read%20offline_data&response_type=code",
			expectedError:    false,
		},
		{
			name:             "Empty input",
			wahooClientId:    "",
			wahooRedirectUri: "",
			expectedResult:   "https://api.wahooligan.com/oauth/authorize?client_id=&redirect_uri=&scope=user_read%20workouts_read%20offline_data&response_type=code",
			expectedError:    false,
		},
		{
			name:             "Invalid redirect URI",
			wahooClientId:    "client456",
			wahooRedirectUri: "invalid_uri",
			expectedResult:   "https://api.wahooligan.com/oauth/authorize?client_id=client456&redirect_uri=invalid_uri&scope=user_read%20workouts_read%20offline_data&response_type=code",
			expectedError:    false,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.GetWahooAuthorizeUrl(tc.wahooClientId, tc.wahooRedirectUri)
			fmt.Println(result)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != nil && result.String() != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, result.String())
			}
		})
	}
}

func TestGetWahooOAuthExchangeURLL(t *testing.T) {

	t.Setenv("WAHOO_CLIENT_ID", "client123")
	t.Setenv("WAHOO_CLIENT_SECRET", "client_secret")
	t.Setenv("REDIRECT_URI", "https://example.com/callback")
	t.Setenv("WAHOO_AUTH_BASE_URL", "https://api.wahooligan.com/oauth/authorize")
	t.Setenv("WAHOO_TOKEN_BASE_URL", "https://api.wahooligan.com/oauth/token")

	testCases := []struct {
		name             string
		wahooClientId    string
		wahooRedirectUri string
		expectedResult   string
		expectedError    bool
	}{
		{
			name:             "Valid input",
			wahooClientId:    "client123",
			wahooRedirectUri: "https://example.com/callback",
			expectedResult:   "https://api.wahooligan.com/oauth/token?client_id=client123&client_secret=https://example.com/callback&code=123&grant_type=authorization_code&redirect_uri=https://example.com/callback",
			expectedError:    false,
		},
		{
			name:             "Empty input",
			wahooClientId:    "",
			wahooRedirectUri: "",
			expectedResult:   "https://api.wahooligan.com/oauth/token?client_id=&client_secret=&code=123&grant_type=authorization_code&redirect_uri=",
			expectedError:    false,
		},
		{
			name:             "Invalid redirect URI",
			wahooClientId:    "client456",
			wahooRedirectUri: "invalid_uri",
			expectedResult:   "https://api.wahooligan.com/oauth/token?client_id=client456&client_secret=invalid_uri&code=123&grant_type=authorization_code&redirect_uri=invalid_uri",
			expectedError:    false,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := utils.GetWahooOAuthExchangeURL(tc.wahooClientId, tc.wahooRedirectUri, "123", tc.wahooRedirectUri)
			fmt.Println(result)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != nil && result.String() != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, result.String())
			}
		})
	}
}
