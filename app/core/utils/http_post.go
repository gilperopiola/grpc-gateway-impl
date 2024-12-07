package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func POST(ctx context.Context, url string, payload any, bearer string, client *http.Client) (int, []byte, error) {

	// Prepare request
	req, err := preparePOST(ctx, url, payload, bearer)
	if err != nil {
		return 0, nil, fmt.Errorf("error preparing POST %s: %w", url, err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return 0, nil, fmt.Errorf("error sending POST %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading successful POST %s: %w", url, err)
	}

	return resp.StatusCode, respBody, nil
}

func preparePOST(ctx context.Context, url string, payload any, bearer string) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	return req, nil
}
