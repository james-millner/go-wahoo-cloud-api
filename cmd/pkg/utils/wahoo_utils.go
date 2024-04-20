package utils

import (
	"bytes"
	"fmt"
	"io"
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
