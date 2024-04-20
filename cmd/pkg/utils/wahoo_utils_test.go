package utils

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/wiremock/go-wiremock"
	"gotest.tools/v3/assert"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

// WIP
func testDownloadingWahooFitFile(t *testing.T) {

	container, network, wiremockPort := startWiremock()
	defer container.Close()
	defer network.Close()

	wiremockClient := wiremock.NewClient("http://localhost:" + wiremockPort)
	defer wiremockClient.Reset()

	fitFileAsBytes := loadTestFitFile()

	_ = wiremockClient.StubFor(wiremock.Get(wiremock.URLPathMatching("/fit.fit")).
		WillReturnResponse(
			wiremock.NewResponse().WithStatus(200).WithHeader("Content-Type", "application/octet-stream").WithBody(string(fitFileAsBytes))))
	defer wiremockClient.Reset()

	reader, err := DownloadFitFileContentsToBuffer("http://localhost:" + wiremockPort + "/fit.fit")
	require.NoError(t, err)

	// Read the actual response body
	actualBody, err := io.ReadAll(reader)
	require.NoError(t, err)

	// Assert that the actual response body matches the expected body
	assert.Equal(t, fitFileAsBytes, actualBody)
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

func loadTestFitFile() []byte {
	file, err := os.ReadFile("testdata/small-fit-file.fit")
	if err != nil {
		fmt.Println("Error reading test fit file:", err)
		return nil
	}
	return file
}
