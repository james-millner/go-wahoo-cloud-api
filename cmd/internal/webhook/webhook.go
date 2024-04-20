package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/james-millner/go-wahoo-cloud-api/cmd/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID int `json:"id"`
}

type File struct {
	URL string `json:"url"`
}

type WorkoutSummary struct {
	ID                  int       `json:"id"`
	AscentAccum         string    `json:"ascent_accum"`
	CadenceAvg          string    `json:"cadence_avg"`
	CaloriesAccum       string    `json:"calories_accum"`
	DistanceAccum       string    `json:"distance_accum"`
	DurationActiveAccum string    `json:"duration_active_accum"`
	DurationPausedAccum string    `json:"duration_paused_accum"`
	DurationTotalAccum  string    `json:"duration_total_accum"`
	HeartRateAvg        string    `json:"heart_rate_avg"`
	PowerBikeNpLast     string    `json:"power_bike_np_last"`
	PowerBikeTssLast    string    `json:"power_bike_tss_last"`
	PowerAvg            string    `json:"power_avg"`
	SpeedAvg            string    `json:"speed_avg"`
	WorkAccum           string    `json:"work_accum"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	File                File      `json:"file"`
	Workout             Workout   `json:"workout"`
}

type Workout struct {
	ID            int       `json:"id"`
	Starts        time.Time `json:"starts"`
	Minutes       int       `json:"minutes"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PlanID        any       `json:"plan_id"`
	WorkoutToken  string    `json:"workout_token"`
	WorkoutTypeID int       `json:"workout_type_id"`
}

type WahooCloudApiResponseBody struct {
	EventType      string         `json:"event_type" validate:"required"`
	WebhookToken   string         `json:"webhook_token" validate:"required"`
	User           User           `json:"user" validate:"required"`
	WorkoutSummary WorkoutSummary `json:"workout_summary" validate:"required"`
}

func Callback() func(w http.ResponseWriter, r *http.Request) {

	log.Println("Callback called")

	return func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Request Body: ", string(requestBody))

		var wahooWorkout WahooCloudApiResponseBody
		jErr := json.Unmarshal(requestBody, &wahooWorkout)
		if jErr != nil {
			fmt.Println("Error unmarshalling JSON:", jErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tokenValidator := validator.New(validator.WithRequiredStructEnabled())
		err = tokenValidator.Struct(wahooWorkout)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", jErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Print the decoded data
		log.Println("Decoded data: ", wahooWorkout)
		rEncErr := enc.Encode(wahooWorkout)
		if rEncErr != nil {
			fmt.Println("Error encoding JSON response:", rEncErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sdkConfig, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Printf("Couldn't load default configuration. Here's why: %v\n", err)
			return
		}

		// Create S3 service client
		svc := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("https://fly.storage.tigris.dev")
			o.Region = "auto"
		})

		// Download the fit file
		reader, err := utils.DownloadFitFileContentsToBuffer(wahooWorkout.WorkoutSummary.File.URL)
		if err != nil {
			fmt.Println("Error downloading file:", err)
			http.Error(w, "Internal Server Error. Unable to download fit file.", http.StatusInternalServerError)
		}

		_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(strconv.Itoa(wahooWorkout.WorkoutSummary.Workout.ID) + ".fit"),
			Body:   reader,
		})
		if err != nil {
			log.Printf("Couldn't upload file. Here's why: %v\n", err)
		}

	}
}
