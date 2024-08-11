package gpt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/apimodels"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Gpt API -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.GptAPI = &GptAPI{}

type GptAPI struct {
	httpClient *http.Client
	postFn     func(context.Context, *http.Client, string, any) (int, []byte, error)
}

func NewAPI(postFn func(context.Context, *http.Client, string, any) (int, []byte, error)) core.GptAPI {
	return &GptAPI{
		httpClient: &http.Client{Timeout: 90},
		postFn:     postFn,
	}
}

/* -~-~-~- Types -~-~-~- */

type Endpoint string

const (
	ChatEndpoint Endpoint = "/chat/completions"
)

/* -~-~-~- Chat Completion Endpoint -~-~-~- */

// NewCompletion generates a text completion from the Gpt API.
func (api *GptAPI) NewCompletion(ctx context.Context, prompt, content string) (string, error) {

	var (
		path = string(ChatEndpoint)
		req  = apimodels.Request{
			Model: apimodels.GPT4oMini,
			Messages: []apimodels.Message{
				{Role: "user", Content: prompt},
				{Role: "user", Content: content},
			},
		}
	)

	status, body, err := api.postFn(ctx, api.httpClient, path, req)
	if err != nil {
		return "", err
	}

	core.LogAPICall(path, status, body)

	var response apimodels.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshalling chat completion response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in completion response")
	}

	return response.Choices[0].Message.Content, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
