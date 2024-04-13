package oauth

import (
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/oauth"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOAuthHappyPath_RedirectCorrectly(t *testing.T) {

	t.Setenv("WAHOO_CLIENT_ID", "client123")
	t.Setenv("WAHOO_CLIENT_SECRET", "client_secret")
	t.Setenv("REDIRECT_URI", "https://example.com/callback")
	t.Setenv("WAHOO_AUTH_BASE_URL", "https://api.wahooligan.com/oauth/authorize")
	t.Setenv("WAHOO_TOKEN_BASE_URL", "https://api.wahooligan.com/oauth/token")

	// Test case for happy path
	request, _ := http.NewRequest("GET", "/authorize", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(oauth.Authorize())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)
	assert.Equal(t,
		"https://api.wahooligan.com/oauth/authorize?client_id=client123&redirect_uri="+
			"https://example.com/callback&scope=user_read%20workouts_read%20offline_data&response_type=code",
		response.Result().Header.Get("Location"))
}
