package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GET(ctx context.Context, url string, urlParams map[string]string, bearer string, client *http.Client) (int, []byte, error) {

	// Prepare request
	req, err := prepareGET(ctx, url, urlParams, bearer)
	if err != nil {
		return 0, nil, fmt.Errorf("error preparing GET %s: %w", url, err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return 0, nil, fmt.Errorf("error sending GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading successful GET %s: %w", url, err)
	}

	return resp.StatusCode, respBody, nil
}

func prepareGET(ctx context.Context, url string, urlParams map[string]string, bearer string) (*http.Request, error) {
	var err error
	if url, err = AddQueryParamsToURL(url, urlParams); err != nil || url == "" {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error preparing GET %s: %w", url, err)
	}

	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	return req, nil
}

func AddQueryParamsToURL(baseURL string, queryParams map[string]string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL %s: %w", baseURL, err)
	}

	urlQuery := parsedURL.Query()
	for key, val := range queryParams {
		urlQuery.Set(key, val)
	}
	baseURL += "?" + urlQuery.Encode()

	return baseURL, nil
}
