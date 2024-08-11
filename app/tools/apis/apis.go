package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/gpt"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/weather"
)

var _ core.APIs = &APIs{}

type APIs struct {
	core.GptAPI
	core.WeatherAPI
}

func NewAPIs() *APIs {

	gptPostFn := makePostFunc("https://api.gpt.com/v1", "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")
	gptAPI := gpt.NewAPI(gptPostFn)

	weatherAPIGetFn := makeGetFunc("https://api.weathermap.org/data/2.5", "")
	weatherAPI := weather.NewAPI(weatherAPIGetFn)

	return &APIs{gptAPI, weatherAPI}
}

type getFn func(ctx context.Context, client *http.Client, url string, urlParams map[string]string) (int, []byte, error)

func makeGetFunc(baseURL, auth string) getFn {
	return func(ctx context.Context, client *http.Client, url string, urlParams map[string]string) (int, []byte, error) {

		// Prepare URL with params.
		url = baseURL + url
		if len(urlParams) > 0 {
			buf := bytes.NewBufferString("?")
			for key, val := range urlParams {
				buf.WriteString(fmt.Sprintf("%s=%s&", key, val))
			}
			url += buf.String()[:buf.Len()-1] // remove trailing '&'
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return 0, nil, fmt.Errorf("error creating GET %s request: %w", url, err)
		}

		req.Header.Set("Content-Type", "application/json")

		if auth != "" {
			req.Header.Set("Authorization", "Bearer "+auth)
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, fmt.Errorf("error sending GET %s: %w", url, err)
		}
		defer resp.Body.Close()

		// We got a response back, read it.
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, nil, fmt.Errorf("success sending GET %s, error reading response: %w", url, err)
		}

		return resp.StatusCode, respBody, nil
	}
}

type postFn func(ctx context.Context, client *http.Client, url string, payload any) (int, []byte, error)

func makePostFunc(baseURL, auth string) postFn {
	return func(ctx context.Context, client *http.Client, url string, payload any) (int, []byte, error) {

		// Prepare request.
		url = baseURL + url

		body, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("error marshalling POST %s payload: %w", url, err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
		if err != nil {
			return 0, nil, fmt.Errorf("error creating POST %s request: %w", url, err)
		}

		req.Header.Set("Content-Type", "application/json")

		if auth != "" {
			req.Header.Set("Authorization", "Bearer "+auth)
		}

		// Send request.
		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, fmt.Errorf("error sending POST %s: %w", url, err)
		}
		defer resp.Body.Close()

		// We got a response back, read it.
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, nil, fmt.Errorf("success sending POST %s, error reading response: %w", url, err)
		}

		return resp.StatusCode, respBody, nil
	}
}
