package health

import (
	"encoding/json"
	"net/http"
)

// Response Struct for a simple health check
type Response struct {
	Message string `json:"message"`
}

// Health endpoint
func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		resp := &Response{
			Message: "OK",
		}

		_ = enc.Encode(resp)
	}
}
