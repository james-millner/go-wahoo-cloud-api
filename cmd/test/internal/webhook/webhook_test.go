package webhook

import (
	"encoding/json"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/internal/webhook"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestWahooCallback_HappyPath(t *testing.T) {

	str := "{\"event_type\":\"workout_summary\",\"webhook_token\":\"b50faa0a-a399-40a7-9c5b-321e9af299df\",\"user\":{\"id\":1120489},\"workout_summary\":{\"id\":252869305,\"ascent_accum\":\"179.0\",\"cadence_avg\":\"67.0\",\"calories_accum\":\"438.0\",\"distance_accum\":\"24323.54\",\"duration_active_accum\":\"3557.0\",\"duration_paused_accum\":\"121.0\",\"duration_total_accum\":\"3678.0\",\"heart_rate_avg\":\"153.0\",\"power_bike_np_last\":\"154.0\",\"power_bike_tss_last\":\"45.0\",\"power_avg\":\"122.0\",\"speed_avg\":\"6.84\",\"work_accum\":\"435276.0\",\"created_at\":\"2024-04-12T18:36:11.000Z\",\"updated_at\":\"2024-04-12T18:36:11.000Z\",\"file\":{\"url\":\"https://cdn.wahooligan.com/wahoo-cloud/production/uploads/workout_file/file/A8dm1z2TPq-mXCZ_5KrKtg/2024-04-12-173445-ELEMNT_BOLT_A6D5-177-0.fit\"},\"workout\":{\"id\":281788767,\"starts\":\"2024-04-12T17:34:45.000Z\",\"minutes\":61,\"name\":\"Cycling\",\"created_at\":\"2024-04-12T18:36:11.000Z\",\"updated_at\":\"2024-04-12T18:36:11.000Z\",\"plan_id\":null,\"workout_token\":\"ELEMNT BOLT A6D5:177\",\"workout_type_id\":0}}}"
	expectedResponseBody := unMarshallResponse(str)

	request, _ := http.NewRequest("POST", "/webhook", strings.NewReader(str))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook.Callback())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %v", response.Code)
	}

	actualResponseBody := unMarshallResponse(response.Body.String())

	if !reflect.DeepEqual(expectedResponseBody, actualResponseBody) {
		t.Errorf("Expected response body: to be:")
		t.Errorf("%v", expectedResponseBody)
		t.Errorf("But received: ")
		t.Errorf("%v", actualResponseBody)
	}
}

func TestWahooCallback_InvalidJson(t *testing.T) {

	str := "{id\":0}}}"

	request, _ := http.NewRequest("POST", "/webhook", strings.NewReader(str))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook.Callback())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, but got %v", response.Code)
	}
}

func TestWahooCallback_InvalidWorkoutSummaryJson(t *testing.T) {

	str := "{\"my_key\":\"my_value\"}"

	request, _ := http.NewRequest("POST", "/webhook", strings.NewReader(str))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook.Callback())

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(response, request)

	actualResponseBody := unMarshallResponse(response.Body.String())

	if actualResponseBody != (webhook.WahooCloudApiResponseBody{}) {
		t.Errorf("Expected response body to be a non-empty struct")
	}
}

func unMarshallResponse(wahooRequestBody string) webhook.WahooCloudApiResponseBody {
	var wahooWorkout webhook.WahooCloudApiResponseBody
	_ = json.Unmarshal([]byte(wahooRequestBody), &wahooWorkout)
	return wahooWorkout
}
