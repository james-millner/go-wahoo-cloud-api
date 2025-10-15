package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func DownloadFitFileContentsToBuffer(wahooFitUrl string) (*bytes.Reader, error) {
	resp, err := http.Get(wahooFitUrl)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Convert response body to io.Reader
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())
	return reader, nil
}

func PostFitFileToExternalService(fileData []byte, fileName string, serviceURL string) error {
	// Create a buffer to write our multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("error creating form file: %w", err)
	}

	// Write the file data to the form
	_, err = part.Write(fileData)
	if err != nil {
		return fmt.Errorf("error writing file data: %w", err)
	}

	// Close the multipart writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", serviceURL, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the content type header with the multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("external service returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully posted FIT file to external service. Status: %d", resp.StatusCode)
	return nil
}
