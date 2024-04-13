package health

import (
	"encoding/json"
	"net/http"
)

// HealthResponse Struct
type HealthResponse struct {
	Message string `json:"message"`
}

// HealthHandler endpoint
func HealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		resp := &HealthResponse{
			Message: "OK",
		}

		enc.Encode(resp)
	}
}
