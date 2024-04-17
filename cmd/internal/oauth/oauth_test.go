package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/ory/dockertest/v3"
	"github.com/wiremock/go-wiremock"
	"io"
	"log"
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

	request, _ := http.NewRequest("GET", "/authorize", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(Authorize())
	handler.ServeHTTP(response, request)

	assert.Equal(t,
		response.Result().Header.Get("Location"),
		"https://api.wahooligan.com/oauth/authorize?client_id=client123&redirect_uri="+
			"https://example.com/callback&scope=user_read%20workouts_read%20offline_data&response_type=code")
}

func TestAuthCallback_AuthCodeReceived_HappyPath(t *testing.T) {

	container, network, wiremockPort := startWiremock()
	defer container.Close()
	defer network.Close()

	wiremockClient := wiremock.NewClient("http://localhost:" + wiremockPort)
	defer wiremockClient.Reset()

	t.Setenv("WAHOO_TOKEN_BASE_URL", "http://localhost:"+wiremockPort+"/oauth/token")

	_ = wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/oauth/token")).
		WillReturnResponse(
			wiremock.NewResponse().WithStatus(200).WithJSONBody(map[string]any{
				"access_token":  "my_access_token",
				"token_type":    "Bearer",
				"expires_in":    1234,
				"refresh_token": "my_refresh_token",
				"scope":         "user_read workouts_read offline_data",
				"created_at":    123123123,
			})))
	defer wiremockClient.Reset()

	request, _ := http.NewRequest("GET", "/?code=abc", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthCallback())
	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 200)

	expectedResponseBody := unMarshallResponse(response.Body.String())
	assert.Equal(t, expectedResponseBody.AccessToken, "my_access_token")
}

func TestAuthCallback_AuthCodeReceived_WahooUnavailable(t *testing.T) {

	container, network, wiremockPort := startWiremock()
	defer container.Close()
	defer network.Close()

	wiremockClient := wiremock.NewClient("http://localhost:" + wiremockPort)
	defer wiremockClient.Reset()

	t.Setenv("WAHOO_TOKEN_BASE_URL", "http://localhost:"+wiremockPort+"/oauth/token")

	_ = wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/oauth/token")).
		WillReturnResponse(
			wiremock.NewResponse().WithStatus(500)))

	request, _ := http.NewRequest("GET", "/?code=abc", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthCallback())
	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 500)
}

func TestAuthCallback_AuthCodeReceived_JsonTokenPayloadChanged(t *testing.T) {

	container, network, wiremockPort := startWiremock()
	defer container.Close()
	defer network.Close()

	wiremockClient := wiremock.NewClient("http://localhost:" + wiremockPort)
	defer wiremockClient.Reset()

	t.Setenv("WAHOO_TOKEN_BASE_URL", "http://localhost:"+wiremockPort+"/oauth/token")

	_ = wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/oauth/token")).
		WillReturnResponse(
			wiremock.NewResponse().WithStatus(200).WithJSONBody(map[string]any{
				"refresh_tokened": "my_refresh_token",
				"scopey":          "user_read workouts_read offline_data",
			})))
	defer wiremockClient.Reset()

	request, _ := http.NewRequest("GET", "/?code=abc", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthCallback())
	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 500)
}

func startWiremock() (*dockertest.Resource, *dockertest.Network, string) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	network, err := pool.CreateNetwork("backend")
	if err != nil {
		log.Fatalf("Could not create Network to docker: %s \n", err)
	}

	r, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "wiremock",
		Repository: "rodolpheche/wiremock",
		Networks:   []*dockertest.Network{network},
	})

	if err != nil {
		fmt.Printf("Could not start wiremock: %v \n", err)
		return r, network, ""
	}

	wiremockPort := r.GetPort("8080/tcp")
	fmt.Println("wiremock - connecting to : ", wiremockPort)
	if err := pool.Retry(func() error {

		resp, err := http.Get("http://localhost:" + wiremockPort + "/__admin")
		if err != nil {
			fmt.Printf("trying to connect to wiremock on localhost:%s, got : %v \n", wiremockPort, err)
			return err
		}

		fmt.Println("status: ", resp.StatusCode)
		rs, _ := io.ReadAll(resp.Body)
		fmt.Printf("RESPONSE: %s \n", rs)
		return nil
	}); err != nil {
		fmt.Printf("Could not connect to wiremock container: %v \n", err)
		return r, network, ""
	}

	return r, network, wiremockPort
}

func unMarshallResponse(wahooRequestBody string) WahooTokenResponse {
	var tokenResponse WahooTokenResponse
	_ = json.Unmarshal([]byte(wahooRequestBody), &tokenResponse)
	return tokenResponse
}
