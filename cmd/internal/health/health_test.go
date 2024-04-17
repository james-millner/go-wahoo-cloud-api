package health

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint_HappyPath(t *testing.T) {

	request, _ := http.NewRequest("GET", "/healthz", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(Health())
	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, 200)

	actualResponse := unMarshallResponse(response.Body.String())
	assert.Equal(t, "OK", actualResponse.Message)
}

func unMarshallResponse(healthResponseStr string) Response {
	var healthResponse Response
	_ = json.Unmarshal([]byte(healthResponseStr), &healthResponse)
	return healthResponse
}
