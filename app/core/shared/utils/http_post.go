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

	// Prepare request.
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("error marshalling POST %s payload: %w", url, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, fmt.Errorf("error creating POST %s request: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	// Send request.
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return 0, nil, fmt.Errorf("error sending POST %s: %w", url, err)
	}
	defer resp.Body.Close()

	// We got a response back, read it.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading successful POST %s response: %w", url, err)
	}

	return resp.StatusCode, respBody, nil
}
